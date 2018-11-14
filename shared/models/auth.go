package models

import jwt "github.com/dgrijalva/jwt-go"

type Exception struct {
	Message string `json:"message"`
}

type TokenAuthentication struct {
	Token string `json:"token" form:"token"`
}

type AuthAdminClaims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

type AuthAppClaims struct {
	AppID string `json:"appId"`
	jwt.StandardClaims
}
