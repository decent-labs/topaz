package crypto

import (
	"encoding/hex"

	multihash "github.com/multiformats/go-multihash"
)

// GetReadableHash ...
func GetReadableHash(digest []byte) (string, error) {
	mhBuf, err := multihash.Encode(digest, multihash.SHA2_256)
	if err != nil {
		return "", err
	}

	mh, err := multihash.Cast(mhBuf)
	if err != nil {
		return "", err
	}

	return mh.B58String(), nil
}

// TransformHashToHex ...
func TransformHashToHex(hash []byte) string {
	return hex.EncodeToString(hash)
}
