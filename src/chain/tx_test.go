package chain

import (
	"bytes"
	"strconv"
	"testing"
)

func createTxInput() InputTx {
	hexString := "029a"
	txid, _ := strconv.ParseUint(hexString, 16, 32)
	txid32 := int32(txid)
	scriptsig := []byte{20, 20, 20, 20}
	return InputTx{
		Txid:      txid32,
		OutInx:    0,
		ScriptSig: scriptsig,
		Sequence:  0,
	}
}

func createTxOutput() OutputTx {
	scriptpubkey := []byte{20, 20, 20, 20}
	return OutputTx{
		Value:        20000,
		ScriptPubKey: scriptpubkey,
	}
}

// TestTxSerDer Test serialisation and deserialisation
func TestTxSerDer(t *testing.T) {
	tx := Tx{
		Version:  10,
		TxInNo:   1,
		TxOutNo:  1,
		Vin:      []InputTx{createTxInput()},
		Vout:     []OutputTx{createTxOutput()},
		LockTime: 10,
	}
	serial := tx.Serialise()

	detx := DeserialiseTx(serial)

	if detx.Version != tx.Version {
		t.Errorf("Version %#v expected: %#v", detx.Version, tx.Version)
	}

	if detx.TxInNo != tx.TxInNo {
		t.Errorf("TxInNo %#v expected: %#v", detx.TxInNo, tx.TxInNo)
	}

	if len(detx.Vin) != len(tx.Vin) {
		t.Errorf("TxInNo %#v expected: %#v", 1, 1)
	}

	if !bytes.Equal(detx.Vin[0].ScriptSig, tx.Vin[0].ScriptSig) {
		t.Errorf("TxInNo %#v expected: %#v", detx.Vin[0].ScriptSig, tx.Vin[0].ScriptSig)
	}

	if !bytes.Equal(detx.Vout[0].ScriptPubKey, tx.Vout[0].ScriptPubKey) {
		t.Errorf("TxInNo %#v expected: %#v", detx.Vout[0].ScriptPubKey, tx.Vout[0].ScriptPubKey)
	}

	if detx.LockTime != tx.LockTime {
		t.Errorf("TxInNo %#v expected: %#v", 1, 1)
	}
}
