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

func NewClient(config *config.Config, httpClient *http.Client) *Client {
	return &Client{
		Lister: NewLister(config, httpClient),
		Leaver: NewLeaver(config, httpClient),
		Joiner: NewJoiner(config, httpClient),
	}
}

type Client struct {
	Lister
	Leaver
	Joiner
}
