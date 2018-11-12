package unit_tests

import (
	"net/http"
	"os"
	"testing"

	"github.com/decentorganization/topaz/api/core/authentication"
	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/api/services/models"
	"github.com/decentorganization/topaz/api/settings"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	TestingT(t)
}

type AuthenticationServicesTestSuite struct{}

var _ = Suite(&AuthenticationServicesTestSuite{})
var t *testing.T

func (s *AuthenticationServicesTestSuite) SetUpSuite(c *C) {
	os.Setenv("GO_ENV", "tests")
	settings.Init()
}

func (suite *AuthenticationServicesTestSuite) TestAdminLogin(c *C) {
	user := models.User{
		Username: "haku",
		Password: "testing",
	}
	response, token := services.AdminLogin(&user)
	assert.Equal(t, http.StatusOK, response)
	assert.NotEmpty(t, token)
}

func (suite *AuthenticationServicesTestSuite) TestAdminLoginIncorrectPassword(c *C) {
	user := models.User{
		Username: "haku",
		Password: "Password",
	}
	response, token := services.AdminLogin(&user)
	assert.Equal(t, http.StatusUnauthorized, response)
	assert.Empty(t, token)
}

func (suite *AuthenticationServicesTestSuite) TestAdminLoginIncorrectUsername(c *C) {
	user := models.User{
		Username: "Username",
		Password: "testing",
	}
	response, token := services.AdminLogin(&user)
	assert.Equal(t, http.StatusUnauthorized, response)
	assert.Empty(t, token)
}

func (suite *AuthenticationServicesTestSuite) TestAdminLoginEmptyCredentials(c *C) {
	user := models.User{
		Username: "",
		Password: "",
	}
	response, token := services.AdminLogin(&user)
	assert.Equal(t, http.StatusUnauthorized, response)
	assert.Empty(t, token)
}

func (suite *AuthenticationServicesTestSuite) TestAdminRefreshToken(c *C) {
	user := models.User{
		Username: "haku",
		Password: "testing",
	}
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenString, err := authBackend.GenerateAdminToken(user.UUID)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	assert.Nil(t, err)

	newToken := services.AdminRefreshToken(token)
	assert.NotEmpty(t, newToken)
}

func (suite *AuthenticationServicesTestSuite) TestAdminRefreshTokenInvalidToken(c *C) {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	newToken := services.AdminRefreshToken(token)
	assert.Empty(t, newToken)
}

func (suite *AuthenticationServicesTestSuite) TestLogout(c *C) {
	user := models.User{
		Username: "haku",
		Password: "testing",
	}
	authBackend := auth.InitJWTAuthenticationBackend()
	tokenString, err := authentication.GenerateAdminToken(user.UUID)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})

	err = services.Logout(tokenString, token)
	assert.Nil(t, err)
}
