package crypto

import (
	"bytes"

	"github.com/cbergoon/merkletree"
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
	return GetReadableHash(b)
}
