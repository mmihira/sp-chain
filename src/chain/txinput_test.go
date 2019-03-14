package chain

import (
	"bytes"
	"testing"
	"spchain/util"
)

// TestTxInputSerDer Test serialisation and deserialisation
func TestTxInputSerDer(t *testing.T) {
	txid32 := util.Init32byteArray(0x01)
	scriptsig := []byte{20, 20, 20, 20}
	txinput := InputTx{
		Txid:      txid32[:],
		OutInx:    0,
		ScriptSig: scriptsig,
		Sequence:  0,
	}

	ser := txinput.Ser()
	txinputdes := DeserialiseInputTx(ser)

	if slen := txinput.ScriptSigLen(); int(slen) != 4 {
		t.Errorf("ScriptSigLen go %d, expect %d", slen, 4)
	}

	if !bytes.Equal(txinput.ScriptSig, scriptsig) {
		t.Errorf("ScriptSig %#v expected: %#v", txinput.ScriptSig, scriptsig)
	}

	if !bytes.Equal(txinput.Txid, txid32[:]) {
		t.Errorf("Txid doesn't equal")
	}

	if txinputdes.OutInx != 0 {
		t.Errorf("OutInx doesn't equal %d", 0)
	}

	if !bytes.Equal(txinputdes.ScriptSig, scriptsig) {
		t.Errorf("ScriptSig %#v expected: %#v", txinputdes.ScriptSig, scriptsig)
	}

	if txinputdes.Sequence != 0 {
		t.Errorf("Sequence got %d, expected %d", txinputdes.Sequence, 0)
	}
}
