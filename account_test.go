package libra_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/go-test/deep"

	libra "github.com/philippgille/libra-sdk-go"
)

const (
	// In the Libra CLI:
	//
	//Decoded: AccountResource {
	// 	authentication_key: 0x8cd377191fe0ef113455c8e8d769f0c0147d5bb618bf195c0af31a05fbfd0969,
	// 	balance: 62500000,
	// 	received_events_count: 1,
	// 	sent_events_count: 4,
	// 	sequence_number: 4,
	// }
	testAcc1StateString = "010000002100000001217da6c6b3e19f1825cfb2676daecce3bf3de03cf26647c78df00b371b25cc9744000000200000008cd377191fe0ef113455c8e8d769f0c0147d5bb618bf195c0af31a05fbfd0969a0acb90300000000010000000000000004000000000000000400000000000000"
	// The account resource part of the above account state blob
	testAcc1ResString = "200000008cd377191fe0ef113455c8e8d769f0c0147d5bb618bf195c0af31a05fbfd0969a0acb90300000000010000000000000004000000000000000400000000000000"
	// The auth key part of the above account resource blob
	testAcc1AuthKey = "8cd377191fe0ef113455c8e8d769f0c0147d5bb618bf195c0af31a05fbfd0969"
)

// TestAccountStateDecoding tests if libra.FromAccountStateBlob(...) works correctly.
// It uses libra.FromAccountResourceBlob(...), so that function must be tested independently.
func TestAccountStateDecoding(t *testing.T) {
	testAccState, err := hex.DecodeString(testAcc1StateString)
	if err != nil {
		t.Fatal(err)
	}
	accState, err := libra.FromAccountStateBlob(testAccState)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(accState.Blob, testAccState) {
		t.Fatal("accState.Blob != expected account state blob")
	}
	accRes := accState.AccountResource
	if diff := deep.Equal(accRes, libra.AccountResource{}); diff == nil {
		t.Fatal("accState.AccountResource has a nil value, but should be an object with values")
	}

	// Construct expected account resource
	expectedAccResBlob, err := hex.DecodeString(testAcc1ResString)
	if err != nil {
		t.Fatal(err)
	}
	expectedAccRes, err := libra.FromAccountResourceBlob(expectedAccResBlob)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(accRes, expectedAccRes); diff != nil {
		t.Fatal(diff)
	}
}

// TestAccountResourceDecoding tests if libra.FromAccountResourceBlob(...) works correctly.
func TestAccountResourceDecoding(t *testing.T) {
	// The string is a substring of the account state tested elsewhere
	testAccRes, err := hex.DecodeString(testAcc1ResString)
	if err != nil {
		t.Fatal(err)
	}
	accRes, err := libra.FromAccountResourceBlob(testAccRes)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(accRes, libra.AccountResource{}); diff == nil {
		t.Fatal("accRes is a nil value, but should be an object with values")
	}
	expectedAuthKey, err := hex.DecodeString(testAcc1AuthKey)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(accRes.AuthKey, expectedAuthKey) {
		t.Fatal("accRes.AuthKey != expected auth key")
	}
	if accRes.Balance != uint64(62500000) {
		t.Fatal("accRes.Balance != 62500000")
	}
	if accRes.ReceivedEvents != uint64(1) {
		t.Fatal("accRes.ReceivedEvents != 1")
	}
	if accRes.SentEvents != uint64(4) {
		t.Fatal("accRes.SentEvents != 4")
	}
	if accRes.SequenceNo != uint64(4) {
		t.Fatal("accRes.SequenceNo != 4")
	}
}
