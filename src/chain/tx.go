package chain

import (
	"bytes"
	"encoding/binary"
)

// Tx An sp-chain transaction
type Tx struct {
	Version  int32
	TxInNo   int64
	TxOutNo  int64
	Vin      []InputTx
	Vout     []OutputTx
	LockTime int32
}

// Serialise Serialise the transaction
func (tx *Tx) Serialise() *bytes.Buffer {
	var buffer bytes.Buffer
	binary.Write(&buffer, littleEndian, tx.Version)
	binary.Write(&buffer, littleEndian, tx.TxInNo)
	for _, tx := range tx.Vin {
		binary.Write(&buffer, littleEndian, tx.Ser().Bytes())
	}
	binary.Write(&buffer, littleEndian, tx.TxOutNo)
	for _, tx := range tx.Vout {
		binary.Write(&buffer, littleEndian, tx.Ser().Bytes())
	}
	binary.Write(&buffer, littleEndian, tx.LockTime)
	return &buffer
}

// DeserialiseTx Deserialize a transaction
func DeserialiseTx(b *bytes.Buffer) Tx {
	var readVersion int32
	binary.Read(b, littleEndian, &readVersion)

	var readTxInNo int64
	binary.Read(b, littleEndian, &readTxInNo)
	inputtxs := []InputTx{}
	for i := int64(0); i < readTxInNo; i++ {
		inputtxs = append(inputtxs, DeserialiseInputTx(b))
	}

	var readTxOutNo int64
	binary.Read(b, littleEndian, &readTxOutNo)
	outputtx := []OutputTx{}
	for i := int64(0); i < readTxOutNo; i++ {
		outputtx = append(outputtx, DeserialiseOutputTx(b))
	}

	var readLockTime int32
	binary.Read(b, littleEndian, &readLockTime)

	return Tx{
		Version:  readVersion,
		TxInNo:   readTxInNo,
		Vin:      inputtxs,
		Vout:     outputtx,
		LockTime: readLockTime,
	}
}

// SerialiseForSign Serialise a transaction for signing
// Restrict to SIGHASH_ALL
func (tx *Tx) SerialiseForSign() *bytes.Buffer {
	var buffer bytes.Buffer
	binary.Write(&buffer, littleEndian, tx.Version)
	binary.Write(&buffer, littleEndian, tx.TxInNo)
	for _, tx := range tx.Vin {
		binary.Write(&buffer, littleEndian, tx.SerSigning())
	}
	binary.Write(&buffer, littleEndian, tx.TxOutNo)
	for _, tx := range tx.Vout {
		binary.Write(&buffer, littleEndian, tx.Ser().Bytes())
	}
	binary.Write(&buffer, littleEndian, tx.LockTime)
	return &buffer
}
