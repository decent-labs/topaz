package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	m "github.com/decentorganization/topaz/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// AuthAdmin is middleware which checks a JWT admin-level token
func AuthAdmin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte(os.Getenv("API_JWT_KEY")), nil
				})
				if err != nil {
					json.NewEncoder(w).Encode(m.Exception{Message: err.Error()})
					return
				}
				if token.Valid {
					// decoded := context.Get(req, "decoded")
					// if decoded["type"] != AuthAdmin {
					// 	json.NewEncoder(w).Encode(Exception{Message: "error unauthorized jwt"})
					// 	return
					// }

					context.Set(req, "decoded", token.Claims)
					next(w, req)
				} else {
					json.NewEncoder(w).Encode(m.Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			json.NewEncoder(w).Encode(m.Exception{Message: "An authorization header is required"})
		}
	})
}

// AuthApp is middleware which checks a JWT admin-level token
func AuthApp(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte(os.Getenv("API_JWT_KEY")), nil
				})
				if err != nil {
					json.NewEncoder(w).Encode(m.Exception{Message: err.Error()})
					return
				}
				if token.Valid {
					context.Set(req, "decoded", token.Claims)
					next(w, req)
				} else {
					json.NewEncoder(w).Encode(m.Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			json.NewEncoder(w).Encode(m.Exception{Message: "An authorization header is required"})
		}
	})
}
