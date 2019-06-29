libra-sdk-go
============

[![GoDoc](https://godoc.org/github.com/philippgille/libra-sdk-go?status.svg)](https://godoc.org/github.com/philippgille/libra-sdk-go) [![Go Report Card](https://goreportcard.com/badge/github.com/philippgille/libra-sdk-go)](https://goreportcard.com/report/github.com/philippgille/libra-sdk-go)

Go SDK for the Libra cryptocurrency

> Note: This is work in progress! The API is not stable and will definitely change in the future!  
> This package uses proper semantic versioning though, so you can use vendoring or Go modules to prevent breaking your build.  
> The changelog of this package can be viewed [here](CHANGELOG.md)

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

    acc := "8cd377191fe0ef113455c8e8d769f0c0147d5bb618bf195c0af31a05fbfd0969"
    accState, err := c.GetAccountState(acc)
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

Develop
-------

The proto files are taken from the [Libra repository](https://github.com/libra/libra), commit `bd8e6dc5434d39bd6b56a3502076353d0787a7ef`.

For updating the `rpc` package you currently need to manually update the proto files, make some changes (e.g. `go_package` option) and then run the Go code generation script: `scripts/generate_rpc.sh`

Related projects
----------------

- Libra SDK for Node.js: https://github.com/perfectmak/libra-core
