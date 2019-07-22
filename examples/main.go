package main

import (
	"fmt"
	"time"

	libra "github.com/philippgille/libra-sdk-go"
)

func main() {
	c, err := libra.NewClient("ac.testnet.libra.org:8000", time.Second)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	acc := "82258967750773ea090d125f9aeb9b9e42f6c7ff6df0742336992bfc52b4096a"
	accState, err := c.GetAccountState(acc)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Raw account state: 0x%x\n", accState.Blob)
	fmt.Println()
	fmt.Printf("Account resource: %v\n", accState.AccountResource)

	txlist, err := c.GetTransactionList(1)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Printf("Transaction list resource: %v\n", txlist.Transactions)
}
