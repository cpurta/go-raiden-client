package address

import (
	"net/http"

	"github.com/cpurta/go-raiden-client/config"
)

var (
	_ Getter = &Client{}
)

// NewClient creates a new address client that provides access to a Raiden node
// ethereum address that is being used.
func NewClient(config *config.Config, httpClient *http.Client) *Client {
	return &Client{
		Getter: NewGetter(config, httpClient),
	}
}

// Client is a address client that allows access to the get address HTTP calls to
// a Raiden node.
type Client struct {
	Getter
}
