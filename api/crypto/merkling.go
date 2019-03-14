package crypto

import (
	"bytes"

	"github.com/cbergoon/merkletree"
	multihash "github.com/multiformats/go-multihash"
)

// MerkleLeaf ...
type MerkleLeaf struct {
	Hash []byte
}

// MerkleLeafs ...
type MerkleLeafs []MerkleLeaf

// CalculateHash ...
func (m MerkleLeaf) CalculateHash() ([]byte, error) {
	return m.Hash, nil
}

// Equals ...
func (m MerkleLeaf) Equals(other merkletree.Content) (bool, error) {
	return bytes.Compare(m.Hash, other.(MerkleLeaf).Hash) == 0, nil
}

// GetMerkleRoot ...
func (ms *MerkleLeafs) GetMerkleRoot() (string, error) {
	var list []merkletree.Content
	for _, obj := range *ms {
		list = append(list, obj)
	}

	t, err := merkletree.NewTree(list)
	if err != nil {
		return "", err
	}

	b := t.MerkleRoot()

	// merkle tree library returns SHA256 hashes
	mhBuf, err := multihash.Encode(b, multihash.SHA2_256)
	if err != nil {
		return "", err
	}

	mh, err := multihash.Cast(mhBuf)
	if err != nil {
		return "", err
	}

	return mh.B58String(), nil
}
