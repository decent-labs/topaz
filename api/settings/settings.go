package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var environments = map[string]string{
	"tests":         "api/settings/tests.json",
	"preproduction": "api/settings/pre.json",
	"production":    "api/settings/prod.json",
}

// Settings ...
type Settings struct {
	PrivateKeyPath     string
	PublicKeyPath      string
	JWTExpirationDelta int
}

var settings = Settings{}
var env = "preproduction"

// Init ...
func Init() {
	env = os.Getenv("GO_ENV")
	if env == "" {
		fmt.Println("Warning: Setting preproduction environment due to lack of GO_ENV value")
		env = "preproduction"
	}
	LoadSettingsByEnv(env)
}

// LoadSettingsByEnv ...
func LoadSettingsByEnv(env string) {
	content, err := ioutil.ReadFile(environments[env])
	if err != nil {
		fmt.Println("Error while reading config file", err)
	}
	settings = Settings{}
	jsonErr := json.Unmarshal(content, &settings)
	if jsonErr != nil {
		fmt.Println("Error while parsing config file", jsonErr)
	}
}

// GetEnvironment ...
func GetEnvironment() string {
	return env
}

// Get ...
func Get() Settings {
	if &settings == nil {
		Init()
	}
	return settings
}

// IsTestEnvironment ...
func IsTestEnvironment() bool {
	return env == "tests"
}
