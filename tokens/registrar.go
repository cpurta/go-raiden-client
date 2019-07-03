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

type registerTokenResponse struct {
	NetworkAddress string `json:"token_network_address"`
}

// Lister is a generic interface to list the Ethereum address associated to the
// Raiden node. It allows for a context to be passed to allow for request timeouts
// and/or deadlines on the response.
type Registrar interface {
	Register(ctx context.Context, tokenAddress common.Address) (common.Address, error)
}

var _ Registrar = &defaultRegistrar{}

// NewLister will return a default address lister for a configured Raiden node.
func NewRegistrar(config *config.Config, httpClient *http.Client) Registrar {
	return &defaultRegistrar{
		baseClient: &util.BaseClient{
			Config:     config,
			HTTPClient: httpClient,
		},
	}
}

type defaultRegistrar struct {
	baseClient *util.BaseClient
}

// List will return the associated Ethereum address to the Raiden node configured
// in the Lister config.
func (lister *defaultRegistrar) Register(ctx context.Context, tokenAddress common.Address) (common.Address, error) {
	var (
		err              error
		registerResponse *registerTokenResponse
		networkAddress   = common.Address{}

		requestURL *url.URL
		request    *http.Request
		response   *http.Response
	)

	if requestURL, err = lister.getRequestURL(tokenAddress); err != nil {
		return networkAddress, err
	}

	if request, err = http.NewRequest("PUT", requestURL.String(), nil); err != nil {
		return networkAddress, err
	}

	request = request.WithContext(ctx)

	if response, err = lister.baseClient.HTTPClient.Do(request); err != nil {
		return networkAddress, err
	}

	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&registerResponse); err != nil {
		return networkAddress, err
	}

	networkAddress = common.HexToAddress(registerResponse.NetworkAddress)

	return networkAddress, nil
}

func (lister *defaultRegistrar) getRequestURL(address common.Address) (*url.URL, error) {
	var (
		err        error
		endpoint   = fmt.Sprintf("%s/api/%s/tokens/%s", lister.baseClient.Config.Host, lister.baseClient.Config.APIVersion, address.Hex())
		requestURL *url.URL
	)

	if requestURL, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	return requestURL, nil
}
