package chain

/*
  Utxo struct is used for serialising and deserialising a utxo.
  Serialised as a json string
*/

type Utxo struct {
	value        int64
	scriptPubKey []byte
}
