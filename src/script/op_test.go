package script

import (
	"encoding/hex"
	"testing"
)

func TestPubKeyV1(t *testing.T) {
	key := []byte{1, 2, 3}

	script := Stack{
		[]Operand{
			PUB_KEY_V1{Key: key},
		},
	}

	stack := Stack{}
	ctxt := ScriptContext{}

	result := true
	for _, s := range script.Contents {
		opResult, err := s.Work(&stack, &ctxt)
		result = result && opResult
		if err != nil {
			t.Errorf(err.Error())
		}
	}

	if _, ok := stack.Top().(PUB_KEY_V1); !ok {
		t.Errorf("Expected type: %#v", PUB_KEY_V1{})
	}

	if !result {
		t.Errorf("Expected false result")
	}
}

func TestOpDup(t *testing.T) {
	key := []byte{1, 2, 3}

	script := Stack{
		[]Operand{
			PUB_KEY_V1{Key: key},
			OP_DUP{},
		},
	}

	stack := Stack{}
	ctxt := ScriptContext{}

	result := true
	for _, s := range script.Contents {
		opResult, err := s.Work(&stack, &ctxt)
		result = result && opResult
		if err != nil {
			t.Errorf(err.Error())
		}
	}

	if _, ok := stack.Top().(PUB_KEY_V1); !ok {
		t.Errorf("Expected type: %#v", PUB_KEY_V1{})
	}

	if _, ok := stack.Second().(PUB_KEY_V1); !ok {
		t.Errorf("Expected type: %#v", PUB_KEY_V1{})
	}

	if !result {
		t.Errorf("Expected false result")
	}
}

func TestOpEqualVerify(t *testing.T) {
	key := []byte{1, 2, 3}

	script := Stack{
		[]Operand{
			PUB_KEY_V1{Key: key},
			OP_DUP{},
			OP_EQUALVERIFY{},
		},
	}

	stack := Stack{}
	ctxt := ScriptContext{}

	result := true
	for _, s := range script.Contents {
		opResult, err := s.Work(&stack, &ctxt)
		result = result && opResult
		if err != nil {
			t.Errorf(err.Error())
		}
	}

	if len(stack.Contents) > 0 {
		t.Errorf("Expected stack to be empty")
	}

	if !result {
		t.Errorf("Expected false result")
	}
}

func TestOpHash160(t *testing.T) {
	hexString := "0250863ad64a87ae8a2fe83c1af1a8403cb53f53e486d8511dad8a04887e5b2352"
	key, _ := hex.DecodeString(hexString)

	script := Stack{
		[]Operand{
			PUB_KEY_V1{Key: key},
			OP_HASH_160{},
		},
	}

	stack := Stack{}
	ctxt := ScriptContext{}

	result := true
	for _, s := range script.Contents {
		opResult, err := s.Work(&stack, &ctxt)
		result = result && opResult
		if err != nil {
			t.Errorf(err.Error())
		}
	}

	if !result {
		t.Errorf("Expected false result")
	}

	if len(stack.Contents) != 1 {
		t.Errorf("Expected 1 items in stack")
	}

	opHashResult := hex.EncodeToString(stack.Top().Data())
	expected := "f54a5851e9372b87810a8e60cdd2e7cfd80b6e31"
	if expected != opHashResult {
		t.Errorf("Expected %s, got %s", expected, opHashResult)
	}
}
