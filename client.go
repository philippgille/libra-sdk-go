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
// The return value is the raw undecoded slice of bytes.
// No proof is validated.
func (c Client) GetAccountState(accountAddr string) ([]byte, error) {
	accountAddrBytes, err := hex.DecodeString(accountAddr)
	if err != nil {
		return nil, err
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
	var knownVersion uint64
	knownVersion = 0 // TODO: Does this make a difference for accounts? Or only for events? Might need to be a method parameter.
	updateLedgerRequest := types.UpdateToLatestLedgerRequest{
		ClientKnownVersion: knownVersion,
		RequestedItems:     requestedItems,
	}
	updateLedgerResponse, err := c.acc.UpdateToLatestLedger(context.Background(), &updateLedgerRequest)
	if err != nil {
		return nil, err
	}

	// We only put one request item in the request, so there should only be one response.
	return updateLedgerResponse.GetResponseItems()[0].GetGetAccountStateResponse().GetAccountStateWithProof().GetBlob().GetBlob(), nil
}

// Close closes the underlying gRPC connection.
func (c Client) Close() {
	c.conn.Close()
}

// NewClient creates a new Libra client.
// It connects to the given validator node via gRPC.
// The connection is kept open until Close() is called on the client.
func NewClient(address string, dialTimeout time.Duration) (Client, error) {
	conn, err := grpc.Dial(address, grpc.WithTimeout(dialTimeout), grpc.WithInsecure())
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
