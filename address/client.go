package address

import (
	"net/http"

	"github.com/cpurta/go-raiden-client/config"
)

var (
	_ Getter = &Client{}
)

func NewClient(config *config.Config, httpClient *http.Client) *Client {
	return &Client{
		Getter: NewGetter(config, httpClient),
	}
}

type Client struct {
	Getter
}
