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

func NewObject(template string) (string, error) {
	return sh.NewObject(template)
}

func PatchLink(root string, path string, childhash string, create bool) (string, error) {
	return sh.PatchLink(root, path, childhash, create)
}

func init() {
	conn := fmt.Sprintf("%s:%s", os.Getenv("IPFS_HOST"), os.Getenv("IPFS_PORT"))
	sh = shell.NewShell(conn)
}
