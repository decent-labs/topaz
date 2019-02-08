package unit_tests

import (
	"os"
	"testing"

	"github.com/decentorganization/topaz/api/api/core/authentication"
	"github.com/decentorganization/topaz/api/api/core/redis"
	"github.com/decentorganization/topaz/api/api/services/shared/models"
	"github.com/decentorganization/topaz/api/api/settings"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	TestingT(t)
}

type AuthenticationBackendTestSuite struct{}

var _ = Suite(&AuthenticationBackendTestSuite{})
var t *testing.T

func (s *AuthenticationBackendTestSuite) SetUpSuite(c *C) {
	os.Setenv("GO_ENV", "tests")
	settings.Init()
}

func (s *AuthenticationBackendTestSuite) TestInitJWTAuthenticationBackend(c *C) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	c.Assert(authBackend, NotNil)
	c.Assert(authBackend.PublicKey, NotNil)
}

func (s *AuthenticationBackendTestSuite) TestGenerateToken(c *C) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenString, err := authBackend.GenerateToken("1234")

	assert.Nil(t, err)
	assert.NotEmpty(t, tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})

	assert.Nil(t, err)
	assert.True(t, token.Valid)
}

func (s *AuthenticationBackendTestSuite) TestAuthenticate(c *C) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	user := &models.User{
		Username: "haku",
		Password: "testing",
	}
	c.Assert(authBackend.Authenticate(user), Equals, true)
}

func (s *AuthenticationBackendTestSuite) TestAuthenticateIncorrectPass(c *C) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	user := models.User{
		Username: "haku",
		Password: "Password",
	}
	c.Assert(authBackend.Authenticate(&user), Equals, false)
}

func (s *AuthenticationBackendTestSuite) TestAuthenticateIncorrectUsername(c *C) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	user := &models.User{
		Username: "Haku",
		Password: "testing",
	}
	c.Assert(authBackend.Authenticate(user), Equals, false)
}

func (s *AuthenticationBackendTestSuite) TestLogout(c *C) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenString, err := authBackend.GenerateToken("1234")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	err = authBackend.Logout(tokenString, token)
	assert.Nil(t, err)

	redisConn := redis.Connect()
	redisValue, err := redisConn.GetValue(tokenString)
	assert.Nil(t, err)
	assert.NotEmpty(t, redisValue)
}

func (s *AuthenticationBackendTestSuite) TestIsInBlacklist(c *C) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenString, err := authBackend.GenerateToken("1234")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	err = authBackend.Logout(tokenString, token)
	assert.Nil(t, err)

	assert.True(t, authBackend.IsInBlacklist(tokenString))
}

func (s *AuthenticationBackendTestSuite) TestIsNotInBlacklist(c *C) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	assert.False(t, authBackend.IsInBlacklist("1234"))
}
