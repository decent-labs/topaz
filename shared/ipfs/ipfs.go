package ipfs

import (
	"bytes"
	"fmt"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
)

var sh *shell.Shell

// Add takes a bytearray and adds it to IPFS
func Add(data []byte) (string, error) {
	return sh.Add(bytes.NewReader(data))
}

// Hash the given data without adding it to IPFS
func Hash(data []byte) (string, error) {
	return sh.Add(bytes.NewReader(data), shell.OnlyHash(true))
}

// NewObject creates a new IPFS object based on the object template provided i.e. "unixfs-dir"
func NewObject(template string) (string, error) {
	return sh.NewObject(template)
}

// PatchLink takes a root hash, a path, childhash, and boolean, returning the resulting root
func PatchLink(root string, path string, childhash string, create bool) (string, error) {
	return sh.PatchLink(root, path, childhash, create)
}

func init() {
	conn := fmt.Sprintf("%s:%s", os.Getenv("IPFS_HOST"), os.Getenv("IPFS_PORT"))
	sh = shell.NewShell(conn)
}
