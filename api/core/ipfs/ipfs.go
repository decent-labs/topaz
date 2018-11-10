package ipfs

import (
	"bytes"
	"fmt"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
)

var sh *shell.Shell

func Add(data []byte) (string, error) {
	return sh.Add(bytes.NewReader(data))
}

func init() {
	shConn := fmt.Sprintf("%s:%s", os.Getenv("IPFS_HOST"), os.Getenv("IPFS_PORT"))
	sh = shell.NewShell(shConn)
}
