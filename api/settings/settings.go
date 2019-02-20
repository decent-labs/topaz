package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var version = "0.1.0"

var environments = map[string]string{
	"tests":         "api/settings/tests.json",
	"preproduction": "api/settings/pre.json",
	"sandbox":       "api/settings/sandbox.json",
	"production":    "api/settings/prod.json",
}

type rootContent struct {
	Version     string `json:"ver"`
	Environment string `json:"env"`
}

// Rc ...
var Rc rootContent

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

	Rc = rootContent{
		Environment: env,
		Version:     version,
	}
	loadSettingsByEnv(env)
}

// loadSettingsByEnv ...
func loadSettingsByEnv(env string) {
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
