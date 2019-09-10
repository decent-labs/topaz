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
// func (ms *MerkleLeafs) GetMerkleRoot() ([]byte, error) {
// 	var list []merkletree.Content
// 	for _, obj := range *ms {
// 		list = append(list, obj)
// 	}

// 	t, err := merkletree.NewTree(list)
// 	if err != nil {
// 		return nil, err
// 	}

// 	b := t.MerkleRoot()

// 	// merkle tree library returns SHA256 hashes
// 	mh, err := multihash.Encode(b, multihash.SHA2_256)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return mh, nil
// }

// RootMaker ...
type RootMaker interface {
	GetRoot([]merkletree.Content) ([]byte, error)
}

// HashEncoder ...
type HashEncoder interface {
	Encode([]byte, uint64) ([]byte, error)
}

// MerkleRootMaker ...
type MerkleRootMaker struct {
	rootMaker   RootMaker
	hashEncoder HashEncoder
}

// NewMerkleRootMaker ...
func NewMerkleRootMaker(rootMaker RootMaker, hashEncoder HashEncoder) MerkleRootMaker {
	return MerkleRootMaker{rootMaker, hashEncoder}
}

// GetRoot ...
func (rm MerkleRootMaker) GetRoot(ms *MerkleLeafs) ([]byte, error) {
	var list []merkletree.Content
	for _, obj := range *ms {
		list = append(list, obj)
	}

	b, err := rm.rootMaker.GetRoot(list)
	if err != nil {
		return nil, err
	}

	// merkle tree library returns SHA256 hashes
	mh, err := rm.hashEncoder.Encode(b, multihash.SHA2_256)
	if err != nil {
		return nil, err
	}

	return mh, nil
}
