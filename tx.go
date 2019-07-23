package libra

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/philippgille/libra-sdk-go/rpc/types"
)

// Transaction is a transaction of Libra Coins.
type Transaction struct {
	RawBytes     []byte
	SenderPubKey []byte
	SenderSig    []byte
}

// TransactionList
type TransactionList struct {
	Transactions []Transaction
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
func (tx Transaction) String() string {
	return fmt.Sprintf("{\"raw_bytes\": \"0x%x\", \"sender_pub_key\": \"0x%x\", \"sender_sig\": \"0x%x\"}",
		tx.RawBytes, tx.SenderPubKey, tx.SenderSig)
}

// FromAccountStateBlob converts an account state blob into an object of the AccountState struct.
func FromTransactionList(transactionList []*types.SignedTransaction) (TransactionList, error) {
	results := TransactionList{}
	for _, x := range transactionList {
		result := Transaction{}
		result.RawBytes = x.GetRawTxnBytes()
		result.SenderPubKey = x.GetSenderPublicKey()
		result.SenderSig = x.GetSenderSignature()
		results.Transactions = append(results.Transactions, result)
	}
	return results, nil
}
