package settings

import "os"

// Rc ...
var Rc rootContent

var version = "0.1.20"

type rootContent struct {
	Version         string          `json:"version"`
	Environment     string          `json:"environment"`
	EthereumContent ethereumContent `json:"ethereum"`
}

type ethereumContent struct {
	EthereumNode string `json:"node"`
}

// GenerateRootContent ...
func GenerateRootContent() {
	Rc = rootContent{
		Version:     version,
		Environment: os.Getenv("GO_ENV"),
		EthereumContent: ethereumContent{
			EthereumNode: os.Getenv("GETH_HOST"),
		},
	}
}
