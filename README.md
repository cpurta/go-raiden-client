# go-raiden-client

[![GoDoc](https://godoc.org/github.com/cpurta/go-raiden-client?status.svg)](https://godoc.org/github.com/cpurta/go-raiden-client)
[![CircleCI](https://circleci.com/gh/cpurta/go-raiden-client?style=svg)](https://circleci.com/gh/cpurta/go-raiden-client)

A client for a Raiden node written in Go. This is meant to be used for backend
development environments where you are able to connect to running Raiden node(s).

## Getting Started

Get the Raiden client within your go path by using `go get github.com/cpurta/go-raiden-client`.

You can then start writing some code using the client. Here is an example that will
connect to a locally running Raiden node and get our Ethereum address that the node
is using.

```go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/cpurta/go-raiden-client/config"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	var (
		err          error
		raidenConfig = &config.Config{
			Host:       "http://localhost:5001",
			APIVersion: "v1",
		}

		raidenClient = NewClient(raidenConfig, http.DefaultClient)
		address common.Address
	)

	// get our token address from the raiden node
	if address, err = raidenClient.Address().Get(context.TODO()); err != nil {
		log.Println("there was an error getting raiden address:", err.Error())
	}

	log.Println("raiden address:", address.Hex())
}
```

## Contributing

If you notice some issues please feel free to create one in the repo with as much
detailed information on how you encountered your issues and any error messages. If
you would like to add on features feel free to fork and submit a pull request.

## LICENSE

Distributed under the [MIT License](./LICENSE)
