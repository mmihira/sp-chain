package key

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

// Key Ecdsa pub/priv key
type Key struct {
	PrivateKey          *btcec.PrivateKey
	PrivateKeyHexString string
	PublicKey           *btcec.PublicKey
	PublicKeyHash       []byte
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

	BtcAddressBytes := append(PubKRipeMDBytes, PubKCheckSum...)

	return Key{
		PrivateKey:          k,
		PrivateKeyHexString: hex.EncodeToString(k.D.Bytes()[:]),
		PublicKey:           (*btcec.PublicKey)(&k.PublicKey),
		PublicKeyHash:       ripe160sha256((*btcec.PublicKey)(&k.PublicKey).SerializeCompressed()),
		BtcAddressBytes:     BtcAddressBytes,
		BtcAddressString:    base58.Encode(BtcAddressBytes),
	}
}
