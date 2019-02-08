package models

import jwt "github.com/dgrijalva/jwt-go"

type AuthKey string

const (
	UserID AuthKey = "userId"
	AppID  AuthKey = "appId"
)

type Exception struct {
	Message string `json:"message"`
}

type TokenAuth struct {
	Token string `json:"token" form:"token"`
}

type AuthClaims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

type AuthAppClaims struct {
	AppID string `json:"appId"`
	jwt.StandardClaims
}
