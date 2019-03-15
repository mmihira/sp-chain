package chain

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func mockBlockHeader() BlockHeader {
	prevBlockHash := "0000000000000000000d9dca531a2a0a179ff72d2fcc9339577c98b369337cff"
	merkleRoot := "5467954125490dad6e04ce8d4eaea7ac5bb1171b49a60d79624d9d6dfc624052"

	prevBlockHashBinary, _ := hex.DecodeString(prevBlockHash)
	merkleRootBinary, _ := hex.DecodeString(merkleRoot)

	return BlockHeader{
		545259520,
		prevBlockHashBinary,
		merkleRootBinary,
		20,
		419668748,
		440532392,
	}
}

// TestBlockHeaderSerialisation Test serialisation and deserialisation
func TestBlockHeaderSerialisation(t *testing.T) {

	header := mockBlockHeader()

	ser := header.Ser()
	dser := DeserialiseBlockHeader(ser)

	if !bytes.Equal(dser.PrevBlockHash, header.PrevBlockHash) {
		t.Errorf("PrevBlockHash %#v expected: %#v", dser.PrevBlockHash, header.PrevBlockHash)
	}

	if !bytes.Equal(dser.MerkleRoot, header.MerkleRoot) {
		t.Errorf("MerkleRoot %#v expected: %#v", dser.MerkleRoot, header.MerkleRoot)
	}

	if dser.Version != header.Version {
		t.Errorf("Version %#v expected: %#v", dser.Version, header.Version)
	}

	if dser.TimeStamp != header.TimeStamp {
		t.Errorf("Timestamp %#v expected: %#v", dser.TimeStamp, header.TimeStamp)
	}

	if dser.Nonce != header.Nonce {
		t.Errorf("Nonce %#v expected: %#v", dser.Nonce, header.Nonce)
	}

	if dser.DifficultyTarget != header.DifficultyTarget {
		t.Errorf("DifficultyTarget %#v expected: %#v", dser.DifficultyTarget, header.DifficultyTarget)
	}

	orgHash := header.Hash()
	newHash := dser.Hash()
	if !bytes.Equal(orgHash[:], newHash[:]) {
		t.Errorf("Block header hashes don't match")
	}
}
