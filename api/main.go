package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var httpClient = http.Client{}

// CreateUserRequest is body of incoming request
type CreateUserRequest struct {
	Name     string
	Email    string
	Password string
}

// CreateTokenRequest is the body of incoming request
type CreateTokenRequest struct {
	Email    string
	Password string
}

// CreateAppRequest is the body of incoming request
type CreateAppRequest struct {
	Interval int
	Name     string
}

// StoreResponse is the response of store request
type StoreResponse struct {
	Hash string
	App  int
}

// CreateUserHandler creates a new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting /create-user handler")

	var ur CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "error executing create new user request: %s" + err.Error()})
		return
	}

	if hp, err := HashPassword(ur.Password); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "error hashing password: %s" + err.Error()})
		return
	}

	u := User{Name: ur.Name, Email: ur.Email, Password: hp}
	if err := db.Create(&u).Error; err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "error saving new user: %s" + err.Error()})
		return
	}

	if err := json.NewEncoder(w).Encode(u); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "error encoding new user: %s" + err.Error()})
		return
	}

	log.Println("finished with /create-user handler")
}

// CreateTokenHandler creates a new JWT for a user
func CreateTokenHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("starting /authenticate (create token) handler")

	var ctr CreateTokenRequest
	if err := json.NewDecoder(req.Body).Decode(&ctr); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "error decoding create token request: %s" + err.Error()})
		return
	}

	var u User
	if err := db.Where("email = ?", ctr.Email).First(&u).Error; err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "error finding user: %s" + err.Error()})
		return
	}

	if match := CheckPasswordHash(ctr.Password, user.Password); match == false {
		json.NewEncoder(w).Encode(Exception{Message: "error authenticating, bad password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Email,
		"password": user.Password,
	})

	if tokenString, err := token.SignedString([]byte(os.Getenv("API_JWT_KEY"))); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "error creating jwt signature: %s", err.Error()})
		return
	}

	if err := json.NewEncoder(w).Encode(JwtToken{Token: tokenString}); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "error encoding jwt: %s", err.Error()})
		return
	}

	log.Println("finished with /authenticate (create token) handler")
}

// CreateAppHandler creates a new app for a user
func CreateAppHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting /create-app handler")

	var ar CreateAppRequest
	if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "error decoding parameters: %s", err.Error()})
		return
	}

	a := App{Interval: ar.Interval, Name: ar.Name}
	if err := db.Create(&a).Error; err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "error saving new app: %s", err.Error()})
		return
	}

	if err := json.NewEncoder(w).Encode(a); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "error encoding new app: %s", err.Error()})
		return
	}

	log.Println("finished with /create-app handler")
}

// StoreHandler takes data and does the thing
func StoreHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting /store handler")

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "error reading /store request body: %s", err.Error()})
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	url := fmt.Sprintf("http://%s:8080", os.Getenv("STORE_HOST"))
	if proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body)); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "error creating proxy /store request to store service: %s", err.Error()})
		return
	}

	if resp, err := httpClient.Do(proxyReq); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "error executing store service request: %s", err.Error()})
		return
	}
	defer resp.Body.Close()

	sr := new(StoreResponse)
	if err = json.NewDecoder(resp.Body).Decode(sr); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "error decoding store service json response: %s", err.Error()})
		return
	}

	if err = json.NewEncoder(w).Encode(sr); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "error encoding store response from api: %s", err.Error()})
		return
	}

	log.Println("finished with /store handler")
}

// VerifyHandler takes data and verifies it's authenticity
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
}

// ReportHandler takes a date range and returns metadata from within it
func ReportHandler(w http.ResponseWriter, r *http.Request) {
}

func main() {
	if i, err := strconv.Atoi(os.Getenv("STARTUP_SLEEP")); err != nil {
		log.Fatalf("missing valid STARTUP_SLEEP environment variable: %s", err.Error())
	}
	time.Sleep(time.Duration(i) * time.Second)

	dbConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", os.Getenv("PQ_HOST"), os.Getenv("PQ_PORT"), os.Getenv("PQ_USER"), os.Getenv("PQ_NAME"))

	if db, err = gorm.Open("postgres", dbConn); err != nil {
		log.Fatalf("couldn't even pretend to open database connection: %s", err.Error())
	}
	defer db.Close()

	http.HandleFunc("/create-user", CreateUserHandler)
	http.HandleFunc("/log-in", CreateTokenHandler)

	http.HandleFunc("/create-app", Authenticate(CreateAppHandler))
	http.HandleFunc("/store", Authenticate(StoreHandler))
	http.HandleFunc("/verify", Authenticate(VerifyHandler))
	http.HandleFunc("/report", Authenticate(ReportHandler))

	log.Println("wake up, api...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
