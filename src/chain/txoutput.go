package chain

import (
	"bytes"
	"encoding/binary"
)

// OutputTx OutputTx
type OutputTx struct {
	Value        int32
	ScriptPubKey []byte
}

// ScriptPubLen the length of the ScriptSig
func (b *OutputTx) ScriptPubLen() byte {
	return byte(len(b.ScriptPubKey))
}

// Ser Serialise the InputTx
func (b *OutputTx) Ser() *bytes.Buffer {
	var ret bytes.Buffer
	binary.Write(&ret, littleEndian, b.Value)
	binary.Write(&ret, littleEndian, b.ScriptPubLen())
	binary.Write(&ret, littleEndian, b.ScriptPubKey)
	return &ret
}

// Deserialize Deserialze bytes to OutputTx
func DeserializeOutputTx(b *bytes.Buffer) OutputTx {
	var readValue int32
	binary.Read(b, littleEndian, &readValue)

	// Read the ScriptPub
	lenToRead := byte(0)
	binary.Read(b, littleEndian, &lenToRead)
	var readPubKey []byte
	for i := 0; i < int(lenToRead); i++ {
		var read byte
		binary.Read(b, littleEndian, &read)
		readPubKey = append(readPubKey, read)
	}

	return OutputTx{
		Value:        readValue,
		ScriptPubKey: readPubKey,
	}
}
