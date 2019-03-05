package chain

import (
	"bytes"
	"strconv"
	"testing"
)

// TestTxInputSerDer Test serialisation and deserialisation
func TestTxInputSerDer(t *testing.T) {
	hexString := "029a"
	txid, _ := strconv.ParseUint(hexString, 16, 32)
	txid32 := int32(txid)
	scriptsig := []byte{20, 20, 20, 20}
	txinput := InputTx{
		Txid:      txid32,
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

	if txinputdes.Txid != txid32 {
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
