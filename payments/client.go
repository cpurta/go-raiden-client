package payments

import (
	"net/http"

	"github.com/cpurta/go-raiden-client/config"
)

var (
	_ Lister    = &Client{}
	_ Initiator = &Client{}
)

func NewClient(config *config.Config, httpClient *http.Client) *Client {
	return &Client{
		Lister:    NewLister(config, httpClient),
		Initiator: NewInitiator(config, httpClient),
	}
}

type Client struct {
	Lister
	Initiator
}
