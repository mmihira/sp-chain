package script

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/ripemd160"
)

/*
The basic script supported is like:
 Unlocking:  SIG PUB_KEY
 Lockign:  OP_DUP OP_HASH160 ab68025513c3dbd2f7b92a94e0581f5d50f654e7 OP_EQUALVERIFY OP_CHECKSIG
*/

var littleEndian = binary.LittleEndian

var (
	OP_EQUALVERIFY_BYTE = byte(0x88)
	OP_DUP_BYTE         = byte(0x76)
	OP_CHECKSIG_BYTE    = byte(0xac)
	SIG_BYTE            = byte(0x01)
	PUB_KEY_V1_BYTE     = byte(0x02)
	PUB_KEY_HASH_BYTE   = byte(0x03)
	OP_HASH_160_BYTE    = byte(0xa9)
)

// sha256 of the byte buffer followed by ripemd160
func ripe160sha256(b []byte) []byte {
	sha := sha256.Sum256(b)
	h := ripemd160.New()
	h.Write(sha[:])
	return h.Sum(nil)
}

// ScriptContext Context which the operands can use when doing
// their operations
type ScriptContext interface {
	SerialiseForSign() *bytes.Buffer
}

// ByteRepresentation Interface for OP_CODES to be represented
// as bytes
type ByteRepresentation interface {
	AsByte() byte
	LenData() byte
	Data() []byte
}

// HasName string names for OP_CODES
type HasName interface {
	Name() string
}

// Operand interface
type Operand interface {
	Work(*Stack, ScriptContext) (bool, error)
	Copy() Operand
	ByteRepresentation
	HasName
}

// OP_DUP Duplicates the top stack item
type OP_DUP struct{}

func (OP_DUP) Work(s *Stack, w ScriptContext) (bool, error) {
	s.DuplicateTop()
	return true, nil
}
func (OP_DUP) AsByte() byte  { return OP_DUP_BYTE }
func (OP_DUP) LenData() byte { return 0x00 }
func (OP_DUP) Data() []byte  { return []byte{0x00} }
func (OP_DUP) Name() string  { return "OP_DUP" }
func (OP_DUP) Copy() Operand {
	return OP_DUP{}
}

// OP_EQUALVERIFY Returns 1 if the inputs are exactly equal, 0 otherwise.
// Then marks transaction as invalid if top stack value is not true.
// The top stack value is removed.
type OP_EQUALVERIFY struct{}

func (OP_EQUALVERIFY) Work(s *Stack, w ScriptContext) (bool, error) {
	first := s.Top()
	second := s.Second()
	firstT, _ := first.(ByteRepresentation)
	secondT, _ := second.(ByteRepresentation)

	if bytes.Equal(firstT.Data(), secondT.Data()) {
		s.PopTwo()
		return true, nil
	} else {
		return false, nil
	}
}
func (OP_EQUALVERIFY) AsByte() byte  { return OP_EQUALVERIFY_BYTE }
func (OP_EQUALVERIFY) LenData() byte { return 0x00 }
func (OP_EQUALVERIFY) Data() []byte  { return []byte{0x00} }
func (OP_EQUALVERIFY) Name() string  { return "OP_EQUALVERIFY" }
func (s OP_EQUALVERIFY) Copy() Operand {
	return OP_EQUALVERIFY{}
}

// OP_CHECKSIG The entire transaction's outputs, inputs, and script
// (from the most recently-executed OP_CODESEPARATOR to the end) are hashed.
// The signature used by OP_OP_CHECKSIG must be a valid signature for this hash and public key.
// If it is, 1 is returned, 0 otherwise.
// We assume all sigs to be chcked as SIGH_HASH_ALL
type OP_CHECKSIG struct{}

func (OP_CHECKSIG) Work(s *Stack, w ScriptContext) (bool, error) {
	txHash := w.SerialiseForSign()

	pubKey := s.Top()
	if _, ok := pubKey.(PUB_KEY_V1); !ok {
		return false, &InvalidType{
			fmt.Sprintf("Invalid type. Expected %s, got %s", PUB_KEY_V1{}.Name(), pubKey.Name()),
		}
	}

	// We expect the pubkKey to be a compressed ecdsa key
	pubKeyParsed, pubKeyParsedError := btcec.ParsePubKey(pubKey.Data(), btcec.S256())
	if pubKeyParsedError != nil {
		return false, &PubKeyParseError{"Failed to parse public key"}
	}

	// Parse the signature
	sig := s.Second()
	if _, ok := sig.(SIG); !ok {
		return false, &InvalidType{
			fmt.Sprintf("Invalid type. Expected %s, got %s", SIG{}.Name(), sig.Name()),
		}
	}
	parsedSig, parseSigError := btcec.ParseDERSignature(sig.Data(), btcec.S256())
	if parseSigError != nil {
		return false, &SigParseError{
			fmt.Sprintf("%s --- %s", "Error parsing signature", parseSigError.Error()),
		}
	}

	// Verify the signature
	verify := parsedSig.Verify(txHash.Bytes(), pubKeyParsed)

	if !verify {
		return false, &SigValidationError{"Signature validation error"}
	}

	return true, nil
}
func (OP_CHECKSIG) AsByte() byte  { return OP_CHECKSIG_BYTE }
func (OP_CHECKSIG) LenData() byte { return 0x00 }
func (OP_CHECKSIG) Data() []byte  { return []byte{0x00} }
func (OP_CHECKSIG) Name() string  { return "OP_CHECKSIG" }
func (s OP_CHECKSIG) Copy() Operand {
	return OP_CHECKSIG{}
}

// SIG a signature value
type SIG struct{ Sig []byte }

func (p SIG) Work(s *Stack, w ScriptContext) (bool, error) {
	s.Push(p.Copy())
	return true, nil
}
func (SIG) AsByte() byte     { return SIG_BYTE }
func (op SIG) LenData() byte { return byte(len(op.Sig)) }
func (op SIG) Data() []byte  { return op.Sig }
func (SIG) Name() string     { return "SIG" }
func (s SIG) Copy() Operand {
	return SIG{Sig: append([]byte{}, s.Sig...)}
}

// PUB_KEY_V1 A version 1 pub_key
type PUB_KEY_V1 struct{ Key []byte }

func (p PUB_KEY_V1) Work(s *Stack, w ScriptContext) (bool, error) {
	s.Push(p.Copy())
	return true, nil
}
func (PUB_KEY_V1) AsByte() byte     { return PUB_KEY_V1_BYTE }
func (op PUB_KEY_V1) LenData() byte { return byte(len(op.Key)) }
func (op PUB_KEY_V1) Data() []byte  { return append([]byte{}, op.Key[:]...) }
func (PUB_KEY_V1) Name() string     { return "PUB_KEY_V1" }
func (op PUB_KEY_V1) Copy() Operand {
	return PUB_KEY_V1{Key: append([]byte{}, op.Key...)}
}

// PUB_K_HASH A version 1 pub_key
type PUB_KEY_HASH struct{ Key []byte }

func (p PUB_KEY_HASH) Work(s *Stack, w ScriptContext) (bool, error) {
	s.Push(p.Copy())
	return true, nil
}
func (PUB_KEY_HASH) AsByte() byte     { return PUB_KEY_HASH_BYTE }
func (op PUB_KEY_HASH) LenData() byte { return byte(len(op.Key)) }
func (op PUB_KEY_HASH) Data() []byte  { return append([]byte{}, op.Key[:]...) }
func (PUB_KEY_HASH) Name() string     { return "PUB_KEY_HASH" }
func (op PUB_KEY_HASH) Copy() Operand {
	return PUB_KEY_HASH{Key: append([]byte{}, op.Key...)}
}

// OP_HASH_160 We use the 25 byte btc address representation
// Take the top element, then do sha256, followed by ripemd160 hash
// We check that it is an operand of type PUB_KEY_V1
type OP_HASH_160 struct{}

func (OP_HASH_160) Work(s *Stack, w ScriptContext) (bool, error) {
	top := s.Top().Data()
	_, ok := s.Top().(PUB_KEY_V1)
	if !ok {
		return false, nil
	}

	hash := ripe160sha256(top)
	s.Pop()
	s.Push(PUB_KEY_HASH{hash})
	return true, nil
}
func (OP_HASH_160) AsByte() byte    { return OP_HASH_160_BYTE }
func (OP_HASH_160) LenData() byte   { return 25 }
func (op OP_HASH_160) Data() []byte { return []byte{0x00} }
func (OP_HASH_160) Name() string    { return "OP_HASH_160" }
func (op OP_HASH_160) Copy() Operand {
	return OP_HASH_160{}
}
