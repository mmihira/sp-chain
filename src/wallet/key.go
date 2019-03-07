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

// Key Ecdsa pub/priv key
type Key struct {
	PrivateKeyInt       *big.Int
	PrivateKeyBytes     []byte
	PrivateKeyHexString string
	PublicKeyX          *big.Int
	PublicKeyXBytes     []byte
	PublicKeyHash 			[]byte
	BtcAddressBytes     []byte
	BtcAddressString    string
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

// sha256 of the byte buffer followed by ripemd160
func ripe160sha256(b []byte) []byte {
	sha := sha256.Sum256(b)
	h := ripemd160.New()
	h.Write(sha[:])
	return h.Sum(nil)
}

// keyFromLibPrivKey Create a Key from a btcec.PrivateKey
func keyFromLibPrivKey(k *btcec.PrivateKey) Key {
	PubKBytes := []byte{0x02}
	// (33 bytes, 1 byte 0x02 (y-coord is even), and 32 bytes corresponding to X coordinate)
	PubKBytes = append(PubKBytes, k.PublicKey.X.Bytes()...)
	PubKRipeMDBytes := ripe160sha256(PubKBytes)
	// Add version byte in front of RIPEMD-160 hash (0x00 for Main Network)
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
		PublicKeyX:          k.PublicKey.X,
		PublicKeyXBytes:     k.PublicKey.X.Bytes(),
		PublicKeyHash: 		   ripe160sha256(k.PublicKey.X.Bytes()),
		BtcAddressBytes:     BtcAddressBytes,
		BtcAddressString:    base58.Encode(BtcAddressBytes),
	}
}
