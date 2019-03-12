package script

import (
	"spchain/chain"
	"spchain/key"
	"testing"
)

// Full test of Pay2PublickKeyHash
func TestP2PKH(t *testing.T) {
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
			OP_DUP{},
			OP_HASH_160{},
			PUB_KEY_V1{key.PublicKeyHash},
			OP_EQUALVERIFY{},
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
