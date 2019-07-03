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

func NewClient(config *config.Config, httpClient *http.Client) *Client {
	return &Client{
		Opener:            NewOpener(config, httpClient),
		Closer:            NewCloser(config, httpClient),
		IncreaseDepositor: NewIncreaseDepositor(config, httpClient),
	}
}

type Client struct {
	Opener
	Closer
	IncreaseDepositor
}
