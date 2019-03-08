package script

import (
	"testing"
	"spchain/key"
)

func TestTop(t *testing.T) {
	stack := Stack {
		[]Operand{OP_DUP{}},
	}

	if _, ok := stack.Top().(OP_DUP) ; !ok {
		t.Errorf("Expected type: %#v", OP_DUP{})
	}
}

func TestSecond(t *testing.T) {
	stack := Stack {
		[]Operand{
			OP_DUP{},
			OP_HASH_160{},
		},
	}

	if _, ok := stack.Second().(OP_DUP) ; !ok {
		t.Errorf("Expected type: %#v", OP_DUP{})
	}

	if _, ok := stack.Top().(OP_HASH_160) ; !ok {
		t.Errorf("Expected type: %#v", OP_HASH_160{})
	}
}

func TestDuplicateTop(t *testing.T) {
	key := []byte{1,2,3}

	stack := Stack {
		[]Operand{
			PUB_KEY_V1{ Key: key },
		},
	}
	stack.DuplicateTop()

	if _, ok := stack.Top().(PUB_KEY_V1) ; !ok {
		t.Errorf("Expected type: %#v", PUB_KEY_V1{})
	}

	if _, ok := stack.Second().(PUB_KEY_V1) ; !ok {
		t.Errorf("Expected type: %#v", PUB_KEY_V1{})
	}

	// Make sure we copy memory bytes instead of just
	// referencing. Change a value and make sure
	// it is not duplicated
	topV, _ :=  stack.Top().(PUB_KEY_V1)
	topV.Key[0] = 99
	sTopV, _ := stack.Second().(PUB_KEY_V1)

	if topV.Key[0] != 99 {
		t.Errorf("Expected %#v, got %#v", 1, topV.Key)
	}

	if sTopV.Key[0] != 1 {
		t.Errorf("Expected %#v, got %#v", 1, sTopV.Key)
	}
}

func TestPush(t *testing.T) {
	stack := Stack {
		[]Operand{ },
	}

	stack.Push(OP_DUP{})

	if _, ok := stack.Top().(OP_DUP) ; !ok {
		t.Errorf("Expected type: %#v", OP_DUP{})
	}
}

func TestPopTwo(t *testing.T) {
	key := []byte{1,2,3}
	stack := Stack {
		[]Operand{
			PUB_KEY_V1{ Key: key },
			OP_DUP{},
			OP_DUP{},
		},
	}

	stack.PopTwo()

	if _, ok := stack.Top().(PUB_KEY_V1) ; !ok {
		t.Errorf("Expected type: %#v", PUB_KEY_V1{})
	}
}

// Test serialization and deserialization
func TestSerDer(t *testing.T) {
	privKeyHexString := "18e14a7b6a307f426a94f8114701e7c8e774e7f9a47e2c2035db29a206321725"
	key := key.ImportFromPrivKeyHexString(privKeyHexString)

	stack := Stack {
		[]Operand{
			PUB_KEY_V1{ key.PublicKey.SerializeCompressed()},
			OP_EQUALVERIFY{},
			OP_HASH_160{},
			OP_DUP{},
			PUB_KEY_V1{ key.PublicKey.SerializeCompressed()},
		},
	}

	ser := stack.Ser()
	newStack := Marshall(ser)

	if _, ok := newStack.Top().(PUB_KEY_V1) ; !ok {
		t.Errorf("Expected type: %#v, got %#v", PUB_KEY_V1{}, newStack.Top())
	}

	if _, ok := newStack.Second().(OP_DUP) ; !ok {
		t.Errorf("Expected type: %#v, got %#v", OP_DUP{}, newStack.Top())
	}
}
