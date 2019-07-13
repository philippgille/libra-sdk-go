package libra

import (
	"context"
	"encoding/hex"
	"time"

	"google.golang.org/grpc"

	"github.com/philippgille/libra-sdk-go/rpc/admission_control"
	"github.com/philippgille/libra-sdk-go/rpc/types"
)

// Client is a libra client that's connected to a validator node via gRPC.
// It allows you to query the account state, send transactions etc.
type Client struct {
	// Only for printing. It was only used to create the gRPC client connection.
	address string
	// Shouldn't need to be used. It was only used to create the AdmissionControlClient.
	conn *grpc.ClientConn
	// Actual client
	acc admission_control.AdmissionControlClient
}

// GetAccountState requests the state of the given account.
func (c Client) GetAccountState(accountAddr string) (AccountState, error) {
	accountAddrBytes, err := hex.DecodeString(accountAddr)
	if err != nil {
		return AccountState{}, err
	}
	// From the generated Go code:
	//
	// Types that are valid to be assigned to RequestedItems:
	//	*RequestItem_GetAccountStateRequest
	//	*RequestItem_GetAccountTransactionBySequenceNumberRequest
	//	*RequestItem_GetEventsByEventAccessPathRequest
	//	*RequestItem_GetTransactionsRequest
	requestedItems := []*types.RequestItem{
		&types.RequestItem{
			RequestedItems: &types.RequestItem_GetAccountStateRequest{
				GetAccountStateRequest: &types.GetAccountStateRequest{
					Address: accountAddrBytes,
				},
			},
		},
	}
	knownVersion := uint64(0) // TODO: Does this make a difference for accounts? Or only for events? Might need to be a method parameter.
	updateLedgerRequest := types.UpdateToLatestLedgerRequest{
		ClientKnownVersion: knownVersion,
		RequestedItems:     requestedItems,
	}
	updateLedgerResponse, err := c.acc.UpdateToLatestLedger(context.Background(), &updateLedgerRequest)
	if err != nil {
		return AccountState{}, err
	}

	// We only put one request item in the request, so there should only be one response.
	accStateBlob := updateLedgerResponse.GetResponseItems()[0].GetGetAccountStateResponse().GetAccountStateWithProof().GetBlob().GetBlob()

	return FromAccountStateBlob(accStateBlob)
}

// GetTransactionList
func (c Client) GetTransactionList() (TransactionList, error) {
	// Types that are valid to be assigned to RequestedItems:
	//	*RequestItem_GetAccountStateRequest
	//	*RequestItem_GetAccountTransactionBySequenceNumberRequest
	//	*RequestItem_GetEventsByEventAccessPathRequest
	//	*RequestItem_GetTransactionsRequest
	requestedItems := []*types.RequestItem{
		&types.RequestItem{
			RequestedItems: &types.RequestItem_GetTransactionsRequest{
				GetTransactionsRequest: &types.GetTransactionsRequest{
					StartVersion: 1,
					Limit:        3,
					FetchEvents:  true,
				},
			},
		},
	}
	knownVersion := uint64(0) // TODO: Does this make a difference for accounts? Or only for events? Might need to be a method parameter.
	updateLedgerRequest := types.UpdateToLatestLedgerRequest{
		ClientKnownVersion: knownVersion,
		RequestedItems:     requestedItems,
	}
	updateLedgerResponse, err := c.acc.UpdateToLatestLedger(context.Background(), &updateLedgerRequest)
	if err != nil {
		return TransactionList{}, err
	}

	txListBlob := updateLedgerResponse.GetResponseItems()[0].GetGetTransactionsResponse().GetTxnListWithProof().GetInfos()

	// just used as a placeholder
	if txListBlob != nil {
	}
	return TransactionList{}, err
}

// SendTx sends a transaction to the connected validator node.
func (c Client) SendTx(tx Transaction) error {
	txRequest := admission_control.SubmitTransactionRequest{
		SignedTxn: &types.SignedTransaction{
			RawTxnBytes:     tx.RawBytes,
			SenderPublicKey: tx.SenderPubKey,
			SenderSignature: tx.SenderSig,
		},
	}
	_, err := c.acc.SubmitTransaction(context.Background(), &txRequest)
	return err
}

// Close closes the underlying gRPC connection.
func (c Client) Close() {
	c.conn.Close()
}

// NewClient creates a new Libra client.
// It connects to the given validator node via gRPC.
// The connection is kept open until Close() is called on the client.
func NewClient(address string, dialTimeout time.Duration) (Client, error) {
	ctxWithTimeout, cancelFunc := context.WithTimeout(context.Background(), dialTimeout)
	defer cancelFunc()
	// We need the grpc.WithBlock() dial option so that the timeout is used for establishing the connection
	// and calling the cancel function via defer doesn't lead to cancelling the connection before the timeout.
	conn, err := grpc.DialContext(ctxWithTimeout, address, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return Client{}, err
	}
	acc := admission_control.NewAdmissionControlClient(conn)
	return Client{
		address: address,
		conn:    conn,
		acc:     acc,
	}, nil
}
