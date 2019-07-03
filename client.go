package raidenclient

import (
	"net/http"

	"github.com/cpurta/go-raiden-client/address"
	"github.com/cpurta/go-raiden-client/channels"
	"github.com/cpurta/go-raiden-client/config"
	"github.com/cpurta/go-raiden-client/connections"
	"github.com/cpurta/go-raiden-client/payments"
	"github.com/cpurta/go-raiden-client/tokens"
)

func NewClient(config *config.Config, httpClient *http.Client) *Client {
	return &Client{
		AddressClient:     address.NewClient(config, httpClient),
		TokensClient:      tokens.NewClient(config, httpClient),
		ChannelsClient:    channels.NewClient(config, httpClient),
		PaymentsClient:    payments.NewClient(config, httpClient),
		ConnectionsClient: connections.NewClient(config, httpClient),
	}
}

type Client struct {
	AddressClient     *address.Client
	TokensClient      *tokens.Client
	ChannelsClient    *channels.Client
	PaymentsClient    *payments.Client
	ConnectionsClient *connections.Client
}

func (client *Client) Address() *address.Client {
	return client.AddressClient
}

func (client *Client) Tokens() *tokens.Client {
	return client.TokensClient
}

func (client *Client) Channels() *channels.Client {
	return client.ChannelsClient
}

func (client *Client) Payments() *payments.Client {
	return client.PaymentsClient
}

func (client *Client) Connections() *connections.Client {
	return client.ConnectionsClient
}
