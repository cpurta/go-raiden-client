package tokens

import (
	"net/http"

	"github.com/cpurta/go-raiden-client/config"
)

var (
	_ Lister        = &Client{}
	_ PartnerLister = &Client{}
	_ Getter        = &Client{}
	_ Registrar     = &Client{}
)

func NewClient(config *config.Config, httpClient *http.Client) *Client {
	return &Client{
		Lister:        NewLister(config, httpClient),
		PartnerLister: NewPartnerLister(config, httpClient),
		Getter:        NewGetter(config, httpClient),
		Registrar:     NewRegistrar(config, httpClient),
	}
}

type Client struct {
	Lister
	PartnerLister
	Getter
	Registrar
}
