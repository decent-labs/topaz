package settings

import "os"

// Rc ...
var Rc rootContent

var version = "0.2.3"

type rootContent struct {
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

// GenerateRootContent ...
func GenerateRootContent() {
	Rc = rootContent{
		Version:     version,
		Environment: os.Getenv("GO_ENV"),
	}
}
