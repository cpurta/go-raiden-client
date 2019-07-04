// Package raidenclient creates a generic client to access sub-clients that
// correspond to the various API calls that any Raiden node supports.
package raidenclient

import (
	"net/http"

	"github.com/cpurta/go-raiden-client/address"
	"github.com/cpurta/go-raiden-client/channels"
	"github.com/cpurta/go-raiden-client/config"
	"github.com/cpurta/go-raiden-client/connections"
	"github.com/cpurta/go-raiden-client/payments"
	"github.com/cpurta/go-raiden-client/pending_transfers"
	"github.com/cpurta/go-raiden-client/tokens"
)

// NewClient will return a Raiden client that is able to access all of the API
// calls that are currently available on a Raiden node. This provides access to
// the various sub-clients that correspond to the various API calls available.
func NewClient(config *config.Config, httpClient *http.Client) *Client {
	return &Client{
		AddressClient:          address.NewClient(config, httpClient),
		TokensClient:           tokens.NewClient(config, httpClient),
		ChannelsClient:         channels.NewClient(config, httpClient),
		PaymentsClient:         payments.NewClient(config, httpClient),
		ConnectionsClient:      connections.NewClient(config, httpClient),
		PendingTransfersClient: pendingtransfers.NewClient(config, httpClient),
	}
}

// Client provides access to API sub-clients that correspond to the various API
// calls that a Raiden node supports.
type Client struct {
	AddressClient          *address.Client
	TokensClient           *tokens.Client
	ChannelsClient         *channels.Client
	PaymentsClient         *payments.Client
	ConnectionsClient      *connections.Client
	PendingTransfersClient *pendingtransfers.Client
}

// Address returns the Address sub-client to access the address being used by the
// Raiden node.
func (client *Client) Address() *address.Client {
	return client.AddressClient
}

// Tokens returns the Tokens sub-client that will be able to register, get, and
// list token networks.
func (client *Client) Tokens() *tokens.Client {
	return client.TokensClient
}

// Channels returns the Channels sub-client that will be able to open, close and
// increase the deposit of a micro-payment channel.
func (client *Client) Channels() *channels.Client {
	return client.ChannelsClient
}

// Payments returns the Payments sub-client that will be able to query all of the
// pending payments.
func (client *Client) Payments() *payments.Client {
	return client.PaymentsClient
}

// Connections returns the Connections sub-client that will be able to list, join
// and leave token networks.
func (client *Client) Connections() *connections.Client {
	return client.ConnectionsClient
}

// PendingTransfers returns the PendingTransfers sub-client that will be able to
// query all pending transfers by token or a channel.
func (client *Client) PendingTransfers() *pendingtransfers.Client {
	return client.PendingTransfersClient
}
