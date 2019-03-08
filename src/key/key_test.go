package key

import (
	"encoding/hex"
	"testing"
)

func TestWallet(t *testing.T) {
	privKeyHexString := "18e14a7b6a307f426a94f8114701e7c8e774e7f9a47e2c2035db29a206321725"
	key := ImportFromPrivKeyHexString(privKeyHexString)

	pubKeyStringBase58 := "1PMycacnJaSqwwJqjawXBErnLsZ7RkXUAs"
	if key.BtcAddressString != pubKeyStringBase58 {
		t.Errorf("Private key does not match")
	}

	pubKeyCompressedShaRipHash := "f54a5851e9372b87810a8e60cdd2e7cfd80b6e31"
	createdPubKeyCompressedShaRipHash := hex.EncodeToString(key.PublicKeyHash)

	// pubKeyCompressedShaRipHash is the hash used in transactions
	// No version byte attached
	if pubKeyCompressedShaRipHash != createdPubKeyCompressedShaRipHash {
		t.Errorf(
			"Expected pubKeyCompressedShaRipHash %s, but got %s",
			pubKeyCompressedShaRipHash,
			createdPubKeyCompressedShaRipHash,
		)
	}
}
