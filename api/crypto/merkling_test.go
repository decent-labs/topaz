package crypto_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/cbergoon/merkletree"
	"github.com/decentorganization/topaz/api/crypto"
)

var oneHash = []byte{1}
var twoHash = []byte{1, 2}

func TestCalculateHash(t *testing.T) {
	leaf := crypto.MerkleLeaf{Hash: oneHash}
	calulatedHash, _ := leaf.CalculateHash()
	if bytes.Compare(calulatedHash, oneHash) != 0 {
		t.Fail()
	}
}

func TestEquals(t *testing.T) {
	one := crypto.MerkleLeaf{Hash: oneHash}
	two := crypto.MerkleLeaf{Hash: twoHash}

	var tests = []struct {
		x, y     crypto.MerkleLeaf
		expected bool
	}{
		{one, one, true},
		{two, two, true},
		{one, two, false},
		{two, one, false},
	}

	for _, tt := range tests {
		if actual, _ := tt.x.Equals(tt.y); actual != tt.expected {
			t.Fail()
		}
	}
}

type mockRootMaker struct {
	t        *testing.T
	expected []merkletree.Content
	b        []byte
	err      error
}

func (rm mockRootMaker) GetRoot(contents []merkletree.Content) ([]byte, error) {
	if len(rm.expected) != len(contents) {
		rm.t.Error("Value contents has the wrong length")
	}
	for i, c := range contents {
		expected, _ := rm.expected[i].CalculateHash()
		actual, _ := c.CalculateHash()
		if string(expected) != string(actual) {
			rm.t.Errorf("Did not get expected contents at %v", i)
		}
	}

	return rm.b, rm.err
}

type mockHashEncoder struct {
	t        *testing.T
	expected []byte
	b        []byte
	err      error
}

func (he mockHashEncoder) Encode(buf []byte, code uint64) ([]byte, error) {
	if string(he.expected) != string(buf) {
		he.t.Error("Wrong value")
	}
	return he.b, he.err
}

var oneLeaf = crypto.MerkleLeaf{Hash: oneHash}
var leaves = crypto.MerkleLeafs{oneLeaf}
var contents = []merkletree.Content{oneLeaf}
var root = []byte{1, 2, 3}

func TestGetRoot_Success(t *testing.T) {
	var encodedRoot = []byte{3, 2, 1}
	rootMaker := mockRootMaker{t, contents, root, nil}
	hashEncoder := mockHashEncoder{t, root, encodedRoot, nil}
	merkleRootMaker := crypto.NewMerkleRootMaker(rootMaker, hashEncoder)
	b, err := merkleRootMaker.GetRoot(&leaves)
	if err != nil {
		t.Error(err.Error())
	} else if string(encodedRoot) != string(b) {
		t.Error("Wrong encoded result")
	}
}

func TestGetRoot_RootError(t *testing.T) {
	errorText := "Can't make root"
	rootMaker := mockRootMaker{t, contents, nil, errors.New(errorText)}
	hashEncoder := mockHashEncoder{t, nil, nil, nil}
	merkleRootMaker := crypto.NewMerkleRootMaker(rootMaker, hashEncoder)
	assertError(t, &merkleRootMaker, errorText)
}

func TestGetRoot_EncodeError(t *testing.T) {
	errorText := "Can't encode"
	rootMaker := mockRootMaker{t, contents, root, nil}
	hashEncoder := mockHashEncoder{t, root, nil, errors.New(errorText)}
	merkleRootMaker := crypto.NewMerkleRootMaker(rootMaker, hashEncoder)
	assertError(t, &merkleRootMaker, errorText)
}

func assertError(t *testing.T, rm *crypto.MerkleRootMaker, text string) {
	_, err := rm.GetRoot(&leaves)
	if err == nil {
		t.Error(text)
	} else if err.Error() != text {
		t.Errorf("Got '%v' instead of '%v'", err.Error(), text)
	}
}
