package libra

// Transaction is a transaction of Libra Coins.
type Transaction struct {
	RawBytes     []byte
	SenderPubKey []byte
	SenderSig    []byte
}
