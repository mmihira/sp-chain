package chain

import (
	"bytes"
	"encoding/binary"
)

var littleEndian = binary.LittleEndian

type BlockHeader struct {
	Version int32
	// Should be 32 bytes
	PrevBlockHash []byte
	// Should be 32 bytes
	MerkleRoot []byte
	TimeStamp int64
	DifficultyTarget int32
	Nonce int32
}

// Ser Serialise the blockheader
func (b BlockHeader) Ser() *bytes.Buffer {
	var ret bytes.Buffer
	binary.Write(&ret, littleEndian, b.Version)
	binary.Write(&ret, littleEndian, b.PrevBlockHash)
	binary.Write(&ret, littleEndian, b.MerkleRoot)
	binary.Write(&ret, littleEndian, b.TimeStamp)
	binary.Write(&ret, littleEndian, b.DifficultyTarget)
	binary.Write(&ret, littleEndian, b.Nonce)
	return &ret
}

// DeserialiseBlockHeader Deserialise the blockheader
func DeserialiseBlockHeader(buff *bytes.Buffer) BlockHeader {
	var ret BlockHeader

	binary.Read(buff, littleEndian, &ret.Version)

	prevBlockHash := [32]byte{}
	binary.Read(buff, littleEndian, &prevBlockHash)

	merkleRoot := [32]byte{}
	binary.Read(buff, littleEndian, &merkleRoot)

	binary.Read(buff, littleEndian, &ret.TimeStamp)
	binary.Read(buff, littleEndian, &ret.DifficultyTarget)
	binary.Read(buff, littleEndian, &ret.Nonce)

	ret.PrevBlockHash = prevBlockHash[:]
	ret.MerkleRoot = merkleRoot[:]

	return ret
}
