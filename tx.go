package libra

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Transaction is a transaction of Libra Coins.
type Transaction struct {
	RawBytes     []byte
	SenderPubKey []byte
	SenderSig    []byte
}

// TransactionList
type TransactionList struct {
	// The whole account state as raw bytes
	Blob                    []byte
	TransactionListResource TransactionListResource
}

// FromTransactionListResourceBlob
func FromTransactionListResourceBlob(transactionListResourceBlob []byte) (TransactionListResource, error) {
	result := TransactionListResource{}

	r := bytes.NewReader(transactionListResourceBlob)

	var balance uint64
	err := binary.Read(r, binary.LittleEndian, &balance)
	if err != nil {
		return result, err
	}
	result.Balance = balance

	return result, nil
}

// TransactionListResource
type TransactionListResource struct {
	AuthKey        []byte
	Balance        uint64
	ReceivedEvents uint64
	SentEvents     uint64
	SequenceNo     uint64
}

// String formats the account state similarly to the Libra CLI.
// Numbers are formatted as string because the numbers are uint64,
// whose max value exceeds JSON's "save integer",
// which can lead to parsing errors.
func (ar TransactionListResource) String() string {
	return fmt.Sprintf("{\"var_1\": \"0x%x\" }",
		ar.AuthKey)
}

// FromAccountStateBlob converts an account state blob into an object of the AccountState struct.
func FromTransactionListBlob(transactionListBlob []byte) (TransactionList, error) {
	result := TransactionList{
		Blob: transactionListBlob,
	}
	return result, nil
}
