package script

import (
	"spchain/chain"
	"spchain/key"
	"strconv"
	"testing"
)

func createTxInput() chain.InputTx {
	hexString := "029a"
	txid, _ := strconv.ParseUint(hexString, 16, 32)
	return chain.InputTx{
		Txid:      int32(txid),
		OutInx:    0,
		ScriptSig: []byte{0},
		Sequence:  0,
	}
}

func createTxOutput() chain.OutputTx {
	return chain.OutputTx{
		Value:        20000,
		ScriptPubKey: []byte{0},
	}
}

func TestCheckSig(t *testing.T) {
	tx := chain.Tx{
		Version:  10,
		TxInNo:   1,
		TxOutNo:  1,
		Vin:      []chain.InputTx{createTxInput()},
		Vout:     []chain.OutputTx{createTxOutput()},
		LockTime: 10,
	}

	privKeyHexString := "18e14a7b6a307f426a94f8114701e7c8e774e7f9a47e2c2035db29a206321725"
	key := key.ImportFromPrivKeyHexString(privKeyHexString)

	txHash := tx.SerialiseForSign().Bytes()
	sig, _ := key.PrivateKey.Sign(txHash)

	script := Stack{
		[]Operand{
			SIG{sig.Serialize()},
			PUB_KEY_V1{key.PublicKey.SerializeCompressed()},
			OP_CHECKSIG{},
		},
	}

	stack := Stack{}

	result := true
	for _, s := range script.Contents {
		opResult, err := s.Work(&stack, &tx)
		result = result && opResult
		if err != nil {
			t.Errorf(err.Error())
		}
	}

	if !result {
		t.Errorf("Expected true result")
	}
}
