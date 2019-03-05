package chain

import (
	"bytes"
	"encoding/binary"
)

var littleEndian = binary.LittleEndian

// InputTx Input Transaction
type InputTx struct {
	Txid      int32
	OutInx    int32
	ScriptSig []byte
	Sequence  int32
}

// ScriptSigLen the length of the ScriptSig
func (b *InputTx) ScriptSigLen() byte {
	return byte(len(b.ScriptSig))
}

// Ser Serialise the InputTx
func (b *InputTx) Ser() *bytes.Buffer {
	var ret bytes.Buffer
	binary.Write(&ret, littleEndian, b.Txid)
	binary.Write(&ret, littleEndian, b.OutInx)
	binary.Write(&ret, littleEndian, b.ScriptSigLen())
	binary.Write(&ret, littleEndian, b.ScriptSig)
	binary.Write(&ret, littleEndian, b.Sequence)
	return &ret
}

// DeserialiseInputTx Deserialze bytes to InputTx
func DeserialiseInputTx(b *bytes.Buffer) InputTx {
	var readTxID int32
	binary.Read(b, littleEndian, &readTxID)

	var readOutInx int32
	binary.Read(b, littleEndian, &readOutInx)

	// Read the ScriptSig
	lenToRead := byte(0)
	binary.Read(b, littleEndian, &lenToRead)
	var readScriptSig []byte
	for i := 0; i < int(lenToRead); i++ {
		var read byte
		binary.Read(b, littleEndian, &read)
		readScriptSig = append(readScriptSig, read)
	}

	// Read the sequence
	var readSequence int32
	binary.Read(b, littleEndian, &readSequence)

	return InputTx{
		Txid:      readTxID,
		OutInx:    readOutInx,
		ScriptSig: readScriptSig,
		Sequence:  readSequence,
	}
}

// SerSigning Serialise the InputTx for signing
// The InputTx should created without a ScriptSig
func (b *InputTx) SerSigning() []byte {
	var ret bytes.Buffer
	binary.Write(&ret, littleEndian, b.Txid)
	binary.Write(&ret, littleEndian, b.OutInx)
	binary.Write(&ret, littleEndian, [1]byte{0x00})
	binary.Write(&ret, littleEndian, b.Sequence)
	return ret.Bytes()
}
