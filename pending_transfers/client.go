package pendingtransfers

import (
	"net/http"

	"github.com/cpurta/go-raiden-client/config"
)

var (
	_ Lister = &Client{}
)

// NewClient allows for all Pending Transfer operations to be performed.
func NewClient(config *config.Config, httpClient *http.Client) *Client {
	return &Client{
		Lister: NewLister(config, httpClient),
	}
}

// Client is a holder for the Pending Transfers lister.
type Client struct {
	Lister
}
