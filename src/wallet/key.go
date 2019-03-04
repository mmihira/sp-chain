package key

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"math/big"
)

type Key struct {
	PrivateKeyInt       *big.Int
	PrivateKeyBytes     []byte
	PrivateKeyHexString string

	PublicKeyX       *big.Int
	PublicKeyXBytes  []byte
	BtcAddressBytes  []byte
	BtcAddressString string
}

// ImportFromPrivKeyHexString Generate a new key from a privatekey string
func ImportFromPrivKeyHexString(s string) Key {
	pbytes, _ := hex.DecodeString(s)
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), pbytes)
	return keyFromLibPrivKey(priv)
}

// NewKey Generate a new key
func NewKey() Key {
	newPrivKey, _ := btcec.NewPrivateKey(btcec.S256())
	return keyFromLibPrivKey(newPrivKey)
}

// keyFromLibPrivKey Create a Key from a btcec.PrivateKey
func keyFromLibPrivKey(k *btcec.PrivateKey) Key {
	PubKBytes := []byte{0x02}
	PubKBytes = append(PubKBytes, k.PublicKey.X.Bytes()...)
	PubKSHA256Bytes := sha256.Sum256(PubKBytes)

	h := ripemd160.New()
	h.Write(PubKSHA256Bytes[:])
	PubKRipeMDBytes := h.Sum(nil)
	PubKRipeMDBytes = append([]byte{0x00}, PubKRipeMDBytes...)
	PubKRipeMDBytesSha256 := sha256.Sum256(PubKRipeMDBytes)
	PubKRipeMDBytesSha256 = sha256.Sum256(PubKRipeMDBytesSha256[:])
	PubKCheckSum := PubKRipeMDBytesSha256[:4]
	fmt.Println(hex.EncodeToString(PubKCheckSum))

	BtcAddressBytes := append(PubKRipeMDBytes, PubKCheckSum...)

	return Key{
		PrivateKeyInt:       k.D,
		PrivateKeyBytes:     k.D.Bytes(),
		PrivateKeyHexString: hex.EncodeToString(k.D.Bytes()[:]),

		PublicKeyX:       k.PublicKey.X,
		PublicKeyXBytes:  k.PublicKey.X.Bytes(),
		BtcAddressBytes:  BtcAddressBytes,
		BtcAddressString: base58.Encode(BtcAddressBytes),
	}
}
