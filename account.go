package libra

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

const (
	// See https://github.com/perfectmak/libra-core/blob/6dc45f7b72aee0162d52a9d53d99cdae1adae6f1/lib/constants/PathValues.ts
	accResourceKey = "01217da6c6b3e19f1825cfb2676daecce3bf3de03cf26647c78df00b371b25cc97"
)

// AccountState represents the state of an account.
type AccountState struct {
	// The whole account state as raw bytes
	Blob []byte
	// The account resource with balance etc.
	AccountResource AccountResource
	// TODO: AccountEvents?
}

// FromAccountStateBlob converts an account state blob into an object of the AccountState struct.
func FromAccountStateBlob(accountStateBlob []byte) (AccountState, error) {
	result := AccountState{
		Blob: accountStateBlob,
	}

	r := bytes.NewReader(accountStateBlob)

	// Inspired by https://github.com/perfectmak/libra-core/blob/6dc45f7b72aee0162d52a9d53d99cdae1adae6f1/lib/client.ts#L250

	var mapEntryCount uint32
	err := binary.Read(r, binary.LittleEndian, &mapEntryCount)
	if err != nil {
		return result, err
	}

	m := make(map[string][]byte, mapEntryCount)

	// For now we only decode the account state.
	// Go through the whole blob to find it (it has a specific map key).
	// TODO: Also take care of events.
	accResourceKeyFound := false
	for mapEntryNo := uint32(1); mapEntryNo <= mapEntryCount; mapEntryNo++ {
		var mapKeyLen uint32
		err := binary.Read(r, binary.LittleEndian, &mapKeyLen)
		if err != nil {
			return result, err
		}
		mapKey := make([]byte, mapKeyLen)
		err = binary.Read(r, binary.LittleEndian, &mapKey)
		if err != nil {
			return result, err
		}

		var mapValLength uint32
		err = binary.Read(r, binary.LittleEndian, &mapValLength)
		if err != nil {
			return result, err
		}
		mapVal := make([]byte, mapValLength)
		err = binary.Read(r, binary.LittleEndian, &mapVal)
		if err != nil {
			return result, err
		}

		mapKeyAsString := hex.EncodeToString(mapKey)
		m[mapKeyAsString] = mapVal

		// TODO: Remove when adding handling of events
		if mapKeyAsString == accResourceKey {
			accResourceKeyFound = true
			break
		}
	}

	if !accResourceKeyFound {
		return result, errors.New("The account state blob didn't contain the data of an account resource or there was an error decoding it")
	}

	accResource, err := FromAccountResourceBlob(m[accResourceKey])
	if err != nil {
		return result, err
	}
	result.AccountResource = accResource

	return result, nil
}

// AccountResource represents an account with its balance etc.
type AccountResource struct {
	AuthKey        []byte
	Balance        uint64
	ReceivedEvents uint64
	SentEvents     uint64
	SequenceNo     uint64
}

// String formats the account state similarly to the Libra CLI.
// Numbers are formatted as string because the numbers are uint64,
// whose max value exceeds JSON's "safe integer",
// which can lead to parsing errors.
//
// Info about JSON's "safe integer":
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Number/MAX_SAFE_INTEGER
func (ar AccountResource) String() string {
	return fmt.Sprintf("{\"authentication_key\": \"0x%x\", \"balance\": \"%d\", \"received_events_count\": \"%d\", \"sent_events_count\": \"%d\", \"sequence_number\": \"%d\"}",
		ar.AuthKey, ar.Balance, ar.ReceivedEvents, ar.SentEvents, ar.SequenceNo)
}

// FromAccountResourceBlob converts an account resource blob into an object of the AccountState struct.
func FromAccountResourceBlob(accountResourceBlob []byte) (AccountResource, error) {
	result := AccountResource{}

	r := bytes.NewReader(accountResourceBlob)

	var authKeyLen uint32
	err := binary.Read(r, binary.LittleEndian, &authKeyLen)
	if err != nil {
		return result, err
	}

	authKey := make([]byte, authKeyLen)
	err = binary.Read(r, binary.LittleEndian, &authKey)
	if err != nil {
		return result, err
	}
	result.AuthKey = authKey

	var balance uint64
	err = binary.Read(r, binary.LittleEndian, &balance)
	if err != nil {
		return result, err
	}
	result.Balance = balance

	var receivedEvents uint64
	err = binary.Read(r, binary.LittleEndian, &receivedEvents)
	if err != nil {
		return result, err
	}
	result.ReceivedEvents = receivedEvents

	var sentEvents uint64
	err = binary.Read(r, binary.LittleEndian, &sentEvents)
	if err != nil {
		return result, err
	}
	result.SentEvents = sentEvents

	var sequenceNo uint64
	err = binary.Read(r, binary.LittleEndian, &sequenceNo)
	if err != nil {
		return result, err
	}
	result.SequenceNo = sequenceNo

	return result, nil
}
