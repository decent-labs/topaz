package main

import (
	"github.com/cbergoon/merkletree"
	"github.com/decentorganization/topaz/api/crypto"
	"github.com/multiformats/go-multihash"
)

type merkletreeRootMaker struct{}

func (rm merkletreeRootMaker) GetRoot(contents []merkletree.Content) ([]byte, error) {
	t, err := merkletree.NewTree(contents)
	if err != nil {
		return nil, err
	}
	return t.MerkleRoot(), nil
}

type multihashEncoder struct{}

func (me multihashEncoder) Encode(buf []byte, code uint64) ([]byte, error) {
	return multihash.Encode(buf, code)
}

var merkleRootMaker = crypto.NewMerkleRootMaker(merkletreeRootMaker{}, multihashEncoder{})

// ProvideMerkleRootMaker ...
func ProvideMerkleRootMaker() crypto.MerkleRootMaker {
	return merkleRootMaker
}
