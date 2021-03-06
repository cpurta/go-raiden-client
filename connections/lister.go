package connections

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/cpurta/go-raiden-client/config"
	"github.com/cpurta/go-raiden-client/util"
	"github.com/ethereum/go-ethereum/common"
)

// Connections represent a map of token network address to a Connection pointer.
type Connections map[common.Address]*Connection

// Lister is an interface to list all token network connections on a given Raiden node.
type Lister interface {
	List(ctx context.Context) (Connections, error)
}

// NewLister will create a default lister that will list out all the token network
// connections for a given Raiden node configuration.
func NewLister(config *config.Config, httpClient *http.Client) Lister {
	return &defaultLister{
		baseClient: &util.BaseClient{
			Config:     config,
			HTTPClient: httpClient,
		},
	}
}

type defaultLister struct {
	baseClient *util.BaseClient
}

// List will list out the list of Connections for the provided Raiden node configuration.
func (lister *defaultLister) List(ctx context.Context) (Connections, error) {
	var (
		err         error
		channels    = make(map[string]*Connection)
		connections = make(map[common.Address]*Connection)

		requestURL *url.URL
		request    *http.Request
		response   *http.Response
	)

	if requestURL, err = lister.getRequestURL(); err != nil {
		return nil, err
	}

	if request, err = http.NewRequest("GET", requestURL.String(), nil); err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)

	if response, err = lister.baseClient.HTTPClient.Do(request); err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&channels); err != nil {
		return nil, err
	}

	for tokenAddress, connection := range channels {
		connections[common.HexToAddress(tokenAddress)] = connection
	}

	return connections, nil
}

func (lister *defaultLister) getRequestURL() (*url.URL, error) {
	var (
		err        error
		endpoint   = fmt.Sprintf("%s/api/%s/connections", lister.baseClient.Config.Host, lister.baseClient.Config.APIVersion)
		requestURL *url.URL
	)

	if requestURL, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	return requestURL, nil
}
