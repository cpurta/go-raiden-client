package tokens

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

// Getter is a generic interface to list the Ethereum address associated to the
// Raiden node. It allows for a context to be passed to allow for request timeouts
// and/or deadlines on the response.
type Getter interface {
	Get(ctx context.Context, address common.Address) (common.Address, error)
}

var _ Getter = &defaultGetter{}

// NewGetter will return a default address Getter for a configured Raiden node.
func NewGetter(config *config.Config, httpClient *http.Client) Getter {
	return &defaultGetter{
		baseClient: &util.BaseClient{
			Config:     config,
			HTTPClient: httpClient,
		},
	}
}

type defaultGetter struct {
	baseClient *util.BaseClient
}

// List will return the associated Ethereum address to the Raiden node configured
// in the Getter config.
func (Getter *defaultGetter) Get(ctx context.Context, tokenAddress common.Address) (common.Address, error) {
	var (
		err            error
		address        string
		networkAddress = common.Address{}

		requestURL *url.URL
		request    *http.Request
		response   *http.Response
	)

	if requestURL, err = Getter.getRequestURL(tokenAddress); err != nil {
		return networkAddress, err
	}

	if request, err = http.NewRequest("GET", requestURL.String(), nil); err != nil {
		return networkAddress, err
	}

	request = request.WithContext(ctx)

	if response, err = Getter.baseClient.HTTPClient.Do(request); err != nil {
		return networkAddress, err
	}

	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&address); err != nil {
		return networkAddress, err
	}

	networkAddress = common.HexToAddress(address)

	return networkAddress, nil
}

func (Getter *defaultGetter) getRequestURL(address common.Address) (*url.URL, error) {
	var (
		err        error
		endpoint   = fmt.Sprintf("%s/api/%s/tokens/%s", Getter.baseClient.Config.Host, Getter.baseClient.Config.APIVersion, address.Hex())
		requestURL *url.URL
	)

	if requestURL, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	return requestURL, nil
}
