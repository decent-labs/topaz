package models

import jwt "github.com/dgrijalva/jwt-go"

type Exception struct {
	Message string `json:"message"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type CreateAdminTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateAppTokenRequest struct {
	AppId int `json:"appId"`
}

type AuthAdminClaims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

type AuthAppClaims struct {
	AppID string `json:"appId"`
	jwt.StandardClaims
}
