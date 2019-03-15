package chain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"spchain/db"
	"spchain/util"
)

// Block An sp-chain transaction
type Block struct {
	Size         int32
	Header       BlockHeader
	TxCount      int64
	Transactions []Tx
}

// Block Serialise a block
func (b *Block) Ser() *bytes.Buffer {
	var ret bytes.Buffer
	binary.Write(&ret, littleEndian, b.Size)
	headerSer := b.Header.Ser()
	binary.Write(&ret, littleEndian, headerSer.Bytes())
	binary.Write(&ret, littleEndian, b.TxCount)
	for _, tx := range b.Transactions {
		binary.Write(&ret, littleEndian, tx.Serialise().Bytes())
	}
	return &ret
}

// Hash The block hash is the double SHA256 hash of the block
func (b *Block) Hash() [32]byte {
	ser := b.Ser().Bytes()
	sha1 := sha256.Sum256(ser)
	sha2 := sha256.Sum256(sha1[:])
	return sha2
}

// HashString The block hash as a hex encoded string
func (b *Block) HashString() string {
	hash := b.Hash()
	return hex.EncodeToString(hash[:])
}

type MerkelResult struct {
	Root [32]byte
	Path [][32]byte
}

// Calculate the merkleRoot for a block
// The last value in the linear merkleRoot representation will
// be the merkleRoot
func (b *Block) CalcMerkle() MerkelResult {
	ret := [][32]byte{}
	workingSet := [][32]byte{}

	// Create array of transaction hashes in order
	for _, tx := range b.Transactions {
		workingSet = append(workingSet, tx.Hash())
	}

	// If not even add 0 hash to the end
	if len(workingSet)%2 != 0 {
		workingSet = append(workingSet, util.Init32byteArray(0x00))
	}

	for len(workingSet) > 1 {
		copiedSet := [][32]byte{}
		for _, hash := range workingSet {
			copiedSet = append(copiedSet, hash)
		}
		workingSet = [][32]byte{}

		for i := 0; i < len(copiedSet); i = +2 {
			hash1 := copiedSet[i]
			hash2 := copiedSet[i+1]
			combined := []byte{}
			combined = append(combined, hash1[:]...)
			combined = append(combined, hash2[:]...)
			combinedHash := sha256.Sum256(combined)
			workingSet = append(workingSet, combinedHash)
			ret = append(ret, combinedHash)
		}
	}
	return MerkelResult{
		Path: ret,
		Root: ret[len(ret)-1],
	}
}

// DeserialiseBlock Deserialise a block
func DeserialiseBlock(buff *bytes.Buffer) Block {
	var ret Block

	binary.Read(buff, littleEndian, &ret.Size)
	ret.Header = DeserialiseBlockHeader(buff)
	binary.Read(buff, littleEndian, &ret.TxCount)
	for txCount := int64(0); txCount < ret.TxCount; txCount++ {
		tx := DeserialiseTx(buff)
		ret.Transactions = append(ret.Transactions, tx)
	}

	return ret
}

// GetBlock
func GetBlock(hash string, db db.Interface) (Block, error) {
	buff, err := db.GetBlock(hash)
	if err != nil {
		return Block{}, err
	}

	return DeserialiseBlock(buff), nil
}

// Save this block into the database
func (b* Block) Save(db db.Interface) error {
	hashStr := b.HashString()
	return db.SaveBlock(hashStr, b.Ser())
}
