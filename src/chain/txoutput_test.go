package chain

import (
	"testing"
	"bytes"
)

// TextTxOutputSerDer Test serialisation and deserialisation
func TestTxOutputSerDer(t *testing.T) {
	scriptpubkey := []byte{20, 20, 20, 20}
	txoutput := OutputTx{
		Value: 20000,
		ScriptPubKey: scriptpubkey,
	}

	ser := txoutput.Ser()
	txoutputdes:= DeserializeOutputTx(ser)

	if slen := txoutputdes.ScriptPubLen(); int(slen) != 4 {
		t.Errorf("ScriptPubLen go %d, expect %d", slen, 4)
	}

	if !bytes.Equal(txoutput.ScriptPubKey, scriptpubkey) {
		t.Errorf("ScriptSig %#v expected: %#v", txoutput.ScriptPubKey, scriptpubkey)
	}
}
