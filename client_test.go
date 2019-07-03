package raidenclient

import (
	"context"
	"log"
	"net/http"

	"github.com/cpurta/go-raiden-client/config"
	"github.com/ethereum/go-ethereum/common"
)

func Example() {
	var (
		err          error
		raidenConfig = &config.Config{
			Host:       "http://localhost:5001",
			APIVersion: "v1",
		}

		raidenClient = NewClient(raidenConfig, http.DefaultClient)
		address      common.Address
	)

	// get our token address from the raiden node

	if address, err = raidenClient.Address().Get(context.TODO()); err != nil {
		log.Println("there was an error getting raiden token address:", err.Error())
	}

	log.Println("raiden token address:", address.Hex())
}
