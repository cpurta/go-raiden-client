package pendingtransfers

import (
	"net/http"

	"github.com/cpurta/go-raiden-client/config"
)

var (
	_ Lister = &Client{}
)

func NewClient(config *config.Config, httpClient *http.Client) *Client {
	return &Client{
		Lister: NewLister(config, httpClient),
	}
}

type Client struct {
	Lister
}
