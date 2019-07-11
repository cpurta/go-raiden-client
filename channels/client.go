package channels

import (
	"net/http"

	"github.com/cpurta/go-raiden-client/config"
)

var (
	_ Opener            = &Client{}
	_ Closer            = &Client{}
	_ IncreaseDepositor = &Client{}
)

// NewClient creates a new client to all channel operations that can be performed
// on a Raiden node. This includes Opening, Closing and Increasing the deposit of
// a channel.
func NewClient(config *config.Config, httpClient *http.Client) *Client {
	return &Client{
		Opener:            NewOpener(config, httpClient),
		Closer:            NewCloser(config, httpClient),
		IncreaseDepositor: NewIncreaseDepositor(config, httpClient),
	}
}

// Client allows for all Channel operations to be perfomed over HTTP calls to a
// Raiden node.
type Client struct {
	Opener
	Closer
	IncreaseDepositor
}
