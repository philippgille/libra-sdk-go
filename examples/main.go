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

	accState, err := c.GetAccountState("8cd377191fe0ef113455c8e8d769f0c0147d5bb618bf195c0af31a05fbfd0969")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Account state: 0x%x", accState)
}
