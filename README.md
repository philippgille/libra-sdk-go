libra-sdk-go
============

Go SDK for the Libra cryptocurrency

> Note: This is work in progress! The API is not stable and will definitely change in the future!

Features
--------

- Get account state (raw bytes)
- Send transaction (raw bytes)

### Roadmap

- Instead of returning the undecoded bytes, a proper account state object will be returned in the future
- Instead of the current `Transaction` struct that only takes `RawBytes`, a higher level transaction struct will be added with fields for the sender and receiver address as well as amount of Libra Coins to send.
- A wallet package will be added that can read a recovery file or take a recovery seed string and create accounts with their public/private keypairs from that
- And much more...

Usage
-----

Example:

```go
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
```

Currently prints:

```
Account state: 0x010000002100000001217da6c6b3e19f1825cfb2676daecce3bf3de03cf26647c78df00b371b25cc9744000000200000008cd377191fe0ef113455c8e8d769f0c0147d5bb618bf195c0af31a05fbfd0969a0acb90300000000010000000000000004000000000000000400000000000000
```
