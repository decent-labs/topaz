package authentication

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"strconv"
	"time"

	"github.com/decentorganization/topaz/api/core/redis"
	"github.com/decentorganization/topaz/api/settings"
	"github.com/decentorganization/topaz/models"
	jwt "github.com/dgrijalva/jwt-go"
)

type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

var (
	tokenDuration, _ = strconv.Atoi(os.Getenv("TOKEN_DURATION_HOURS"))
	expireOffset, _  = strconv.Atoi(os.Getenv("TOKEN_EXPIRE_OFFSET"))
)

var authBackendInstance *JWTAuthenticationBackend

func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}

	return authBackendInstance
}

func (backend *JWTAuthenticationBackend) GenerateAdminToken(userID string) (string, error) {
	claims := models.AuthAdminClaims{
		UserID:         userID,
		StandardClaims: generateStandardClaims(userID),
	}

	return backend.generateToken(claims)
}

func (backend *JWTAuthenticationBackend) GenerateAppToken(appID string) (string, error) {
	claims := models.AuthAppClaims{
		AppID:          appID,
		StandardClaims: generateStandardClaims(appID),
	}

	return backend.generateToken(claims)
}

func generateStandardClaims(resourceID string) jwt.StandardClaims {
	return jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   resourceID,
	}
}

func (backend *JWTAuthenticationBackend) generateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (backend *JWTAuthenticationBackend) Authenticate(suppliedPassword string, dbHash string) bool {
	valid := CheckPasswordHash(suppliedPassword, dbHash)
	return valid
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

func (backend *JWTAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
	redisConn := redis.Connect()
	return redisConn.SetValue(tokenString, tokenString, backend.getTokenRemainingValidity(token.Claims.(jwt.MapClaims)["exp"]))
}

func (backend *JWTAuthenticationBackend) IsInBlacklist(token string) bool {
	redisConn := redis.Connect()
	redisToken, _ := redisConn.GetValue(token)

	if redisToken == nil {
		return false
	}

	return true
}

func getPrivateKey() *rsa.PrivateKey {
	privateKeyFile, err := os.Open(settings.Get().PrivateKeyPath)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {
	publicKeyFile, err := os.Open(settings.Get().PublicKeyPath)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	publicKeyFile.Close()

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
