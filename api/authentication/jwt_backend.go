package authentication

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/decentorganization/topaz/shared/models"
	"github.com/decentorganization/topaz/shared/redis"
	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
)

// JWTAuthenticationBackend ...
type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

var (
	expireOffset, _  = strconv.Atoi(os.Getenv("JWT_EXPIRE_OFFSET"))
	tokenDuration, _ = strconv.Atoi(os.Getenv("JWT_TOKEN_DURATION"))
	pkString         = os.Getenv("JWT_PRIVATE_KEY")
	pubString        = os.Getenv("JWT_PUBLIC_KEY")
)

var authBackendInstance *JWTAuthenticationBackend

// InitJWTAuthenticationBackend ...
func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			publicKey:  getPublicKey(),
		}
	}

	return authBackendInstance
}

// GenerateToken ...
func (backend *JWTAuthenticationBackend) GenerateToken(userID string) (string, error) {
	claims := models.AuthClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(tokenDuration)).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Authenticate ...
func (backend *JWTAuthenticationBackend) Authenticate(suppliedPassword string, dbHash string) bool {
	valid := CheckPasswordHash(suppliedPassword, dbHash)
	return valid
}

// GetToken ...
func (backend *JWTAuthenticationBackend) GetToken(req *http.Request) (*jwt.Token, error) {
	token, err := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return backend.publicKey, nil
	})

	return token, err
}

func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds()) + expireOffset
		}
	}
	return expireOffset
}

// Logout ...
func (backend *JWTAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
	return redis.SetValue(tokenString, tokenString, backend.getTokenRemainingValidity(token.Claims.(jwt.MapClaims)["exp"]))
}

// IsInBlacklist ...
func (backend *JWTAuthenticationBackend) IsInBlacklist(token string) bool {
	redisToken, _ := redis.GetString(token)

	if redisToken == "" {
		return false
	}

	return true
}

func getPrivateKey() *rsa.PrivateKey {
	data, _ := pem.Decode([]byte(pkString))
	if data == nil {
		panic("Can't decode JWT_PRIVATE_KEY")
	}

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {
	data, _ := pem.Decode([]byte(pubString))
	if data == nil {
		panic("Can't decode JWT_PUBLIC_KEY")
	}

	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		panic(err)
	}

	return rsaPub
}
