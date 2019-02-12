package models

import jwt "github.com/dgrijalva/jwt-go"

// AuthKey ...
type AuthKey string

// AuthUser ...
const AuthUser AuthKey = "userId"

// Exception ...
type Exception struct {
	Message string `json:"message"`
}

// TokenAuth ...
type TokenAuth struct {
	Token string `json:"token" form:"token"`
}

// AuthClaims ...
type AuthClaims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}
