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
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/decentorganization/topaz/api/auth"
	m "github.com/decentorganization/topaz/models"
)

var db *gorm.DB
var httpClient = http.Client{}

// CreateUserHandler creates a new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting /create-user handler")

	var ur m.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error executing create new user request: %s" + err.Error()})
		return
	}

	hp, err := auth.HashPassword(ur.Password)
	if err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error hashing password: %s" + err.Error()})
		return
	}

	u := m.User{Name: ur.Name, Email: ur.Email, Password: hp}
	if err := db.Create(&u).Error; err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error saving new user: %s" + err.Error()})
		return
	}

	if err := json.NewEncoder(w).Encode(u); err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error encoding new user: %s" + err.Error()})
		return
	}

	log.Println("finished with /create-user handler")
}

// CreateAdminTokenHandler creates a new JWT admin for a user
func CreateAdminTokenHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("starting /auth-admin (create admin token) handler")

	var ctr m.CreateAdminTokenRequest
	if err := json.NewDecoder(req.Body).Decode(&ctr); err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error decoding create token request: %s" + err.Error()})
		return
	}

	var u m.User
	if err := db.Where("email = ?", ctr.Email).First(&u).Error; err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error finding user: %s" + err.Error()})
		return
	}

	if match := auth.CheckPasswordHash(ctr.Password, u.Password); match == false {
		json.NewEncoder(w).Encode(m.Exception{Message: "error authenticating, bad password"})
		return
	}

	claims := m.AuthAdminClaims{
		string(u.ID),
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "topaz-test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("API_JWT_KEY")))
	if err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error creating jwt signature: %s" + err.Error()})
		return
	}

	if err := json.NewEncoder(w).Encode(m.JwtToken{Token: tokenString}); err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error encoding jwt: %s" + err.Error()})
		return
	}

	log.Println("finished with /auth-admin (create admin token) handler")
}

// CreateAppHandler creates a new app for a user
func CreateAppHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting /create-app handler")

	// decoded := context.Get(r, "decoded")
	// if decoded.Get("type") != AuthAdmin {
	// 	json.NewEncoder(w).Encode(m.Exception{Message: "error unauthorized jwt"})
	// 	return
	// }

	// var ar m.CreateAppRequest
	// if err := json.NewDecoder(r.Body).Decode(&ar); err != nil {
	// 	json.NewEncoder(w).Encode(m.Exception{Message: "error decoding parameters: %s" + err.Error()})
	// 	return
	// }

	// a := m.App{Interval: ar.Interval, Name: ar.Name, UserID: decoded["uid"]}
	// if err := db.Create(&a).Error; err != nil {
	// 	json.NewEncoder(w).Encode(m.Exception{Message: "error saving new app: %s" + err.Error()})
	// 	return
	// }

	// if err := json.NewEncoder(w).Encode(a); err != nil {
	// 	json.NewEncoder(w).Encode(m.Exception{Message: "error encoding new app: %s" + err.Error()})
	// 	return
	// }

	log.Println("finished with /create-app handler")
}

// CreateAppTokenHandler creates a new JWT login for a user
func CreateAppTokenHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("starting /create-app (create app token) handler")

	var ctr m.CreateAppTokenRequest
	if err := json.NewDecoder(req.Body).Decode(&ctr); err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error decoding create token request: %s" + err.Error()})
		return
	}

	var a m.App
	if err := db.Where("id = ?", ctr.AppId).First(&a).Error; err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error finding app: %s" + err.Error()})
		return
	}

	claims := m.AuthAppClaims{
		string(a.ID),
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "topaz-test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("API_JWT_KEY")))
	if err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error creating jwt signature: %s" + err.Error()})
		return
	}

	if err := json.NewEncoder(w).Encode(m.JwtToken{Token: tokenString}); err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error encoding jwt: %s" + err.Error()})
		return
	}

	log.Println("finished with /create-app (create app token) handler")
}

// StoreHandler takes data and does the thing
func StoreHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting /store handler")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error reading /store request body: %s" + err.Error()})
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	url := fmt.Sprintf("http://%s:8080", os.Getenv("STORE_HOST"))
	proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
	if err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error creating proxy /store request to store service: %s" + err.Error()})
		return
	}

	resp, err := httpClient.Do(proxyReq)
	if err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error executing store service request: %s" + err.Error()})
		return
	}
	defer resp.Body.Close()

	sr := new(m.StoreResponse)
	if err := json.NewDecoder(resp.Body).Decode(sr); err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error decoding store service json response: %s" + err.Error()})
		return
	}

	if err := json.NewEncoder(w).Encode(sr); err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error encoding store response from api: %s" + err.Error()})
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
	i, err := strconv.Atoi(os.Getenv("STARTUP_SLEEP"))
	if err != nil {
		log.Fatalf("missing valid STARTUP_SLEEP environment variable: %s", err.Error())
	}
	time.Sleep(time.Duration(i) * time.Second)

	dbConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", os.Getenv("PQ_HOST"), os.Getenv("PQ_PORT"), os.Getenv("PQ_USER"), os.Getenv("PQ_NAME"))

	db, err := gorm.Open("postgres", dbConn)
	if err != nil {
		log.Fatalf("couldn't even pretend to open database connection: %s", err.Error())
	}
	defer db.Close()

	http.HandleFunc("/create-admin", CreateUserHandler)
	http.HandleFunc("/auth-admin", CreateAdminTokenHandler)

	http.HandleFunc("/create-app", auth.AuthAdmin(CreateAppHandler))
	http.HandleFunc("/auth-app", auth.AuthAdmin(CreateAppTokenHandler))

	http.HandleFunc("/store", auth.AuthApp(StoreHandler))
	http.HandleFunc("/verify", auth.AuthApp(VerifyHandler))
	http.HandleFunc("/report", auth.AuthApp(ReportHandler))

	log.Println("wake up, api...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
