package connections

import (
	"net/http"

	"github.com/cpurta/go-raiden-client/config"
)

var (
	_ Lister = &Client{}
	_ Leaver = &Client{}
	_ Joiner = &Client{}
)

// NewClient creates a new Connections client that will be able to List all Open
// connections, Join and Leave connections for a Raiden node.
func NewClient(config *config.Config, httpClient *http.Client) *Client {
	return &Client{
		Lister: NewLister(config, httpClient),
		Leaver: NewLeaver(config, httpClient),
		Joiner: NewJoiner(config, httpClient),
	}
}

// Client allows for list, leave and join operations for a Raiden node.
type Client struct {
	Lister
	Leaver
	Joiner
}
