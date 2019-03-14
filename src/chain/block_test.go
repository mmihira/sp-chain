package chain

import (
	"spchain/key"
	"spchain/script"
	"testing"
	"spchain/util"
	"github.com/davecgh/go-spew/spew"
)

func createTxOutputBlockTest() OutputTx {
	privKeyHexString := "18e14a7b6a307f426a94f8114701e7c8e774e7f9a47e2c2035db29a206321725"
	key := key.ImportFromPrivKeyHexString(privKeyHexString)
	script := script.Stack{
		Contents: []script.Operand{
			script.OP_DUP{},
			script.OP_HASH_160{},
			script.PUB_KEY_V1{Key: key.PublicKeyHash},
			script.OP_EQUALVERIFY{},
			script.OP_CHECKSIG{},
		},
	}

	scriptSer := script.Ser().Bytes()

	return OutputTx{
		Value:        20000,
		ScriptPubKey: scriptSer,
	}
}

func createTxInputBlockTestNoSig() InputTx {
	txid32 := util.Init32byteArray(0x01)
	return InputTx{
		Txid:      txid32[:],
		OutInx:    1,
		ScriptSig: []byte{},
		Sequence:  20,
	}
}

func createTxBlockTest() Tx {
	// Create transaction first with no input scriptsigs
	tx := Tx{
		Version:  10,
		TxInNo:   1,
		TxOutNo:  1,
		Vin:      []InputTx{createTxInputBlockTestNoSig()},
		Vout:     []OutputTx{createTxOutputBlockTest()},
		LockTime: 10,
	}

	//populate the script sigs
	privKeyHexString := "18e14a7b6a307f426a94f8114701e7c8e774e7f9a47e2c2035db29a206321725"
	key := key.ImportFromPrivKeyHexString(privKeyHexString)

	txHash := tx.SerialiseForSign().Bytes()
	sig, _ := key.PrivateKey.Sign(txHash)

	script := script.Stack{
		Contents: []script.Operand{
			script.SIG{Sig: sig.Serialize()},
			script.PUB_KEY_V1{Key: key.PublicKey.SerializeCompressed()},
		},
	}
	scriptSer := script.Ser().Bytes()

	tx.Vin[0].ScriptSig = scriptSer

	return tx
}

func TestMerkle(t *testing.T) {
	block := Block {
		Size: 30,
		Header: mockBlockHeader(),
		TxCount: 1,
		Transactions: []Tx{createTxBlockTest()},
	}

	merkle := block.CalcMerkle()
	spew.Dump(merkle)
	t.Errorf("Size %#v", merkle)
}

// TestBlockSerialisation Test serialisation and deserialisation
func TestBlockSerialisation(t *testing.T) {
	block := Block {
		Size: 30,
		Header: mockBlockHeader(),
		TxCount: 1,
		Transactions: []Tx{createTxBlockTest()},
	}

	ser := block.Ser()
	dser := DeserialiseBlock(ser)

	if dser.Size != block.Size {
		t.Errorf("Size %#v expected: %#v", dser.Size, block.Size)
	}

	if dser.TxCount != block.TxCount {
		t.Errorf("TxCount %#v expected: %#v", dser.TxCount, block.TxCount)
	}
}
