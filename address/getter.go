package address

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

type addressResponse struct {
	OurAddress string `json:"our_address"`
}

// Getter is a generic interface to list the Ethereum address associated to the
// Raiden node. It allows for a context to be passed to allow for request timeouts
// and/or deadlines on the response.
type Getter interface {
	Get(ctx context.Context) (common.Address, error)
}

var _ Getter = &defaultGetter{}

// NewGetter will return a default address lister for a configured Raiden node.
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

// Get will return the associated Ethereum address to the Raiden node configured
// in the Lister config.
func (lister *defaultGetter) Get(ctx context.Context) (common.Address, error) {
	var (
		err     error
		address = common.Address{}

		requestURL      *url.URL
		request         *http.Request
		response        *http.Response
		addressResponse = &addressResponse{}
	)

	if requestURL, err = lister.getRequestURL(); err != nil {
		return address, err
	}

	if request, err = http.NewRequest("GET", requestURL.String(), nil); err != nil {
		return address, err
	}

	request = request.WithContext(ctx)

	if response, err = lister.baseClient.HTTPClient.Do(request); err != nil {
		return address, err
	}

	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&addressResponse); err != nil {
		return address, err
	}

	address = common.HexToAddress(addressResponse.OurAddress)

	return address, nil
}

func (lister *defaultGetter) getRequestURL() (*url.URL, error) {
	var (
		err        error
		endpoint   = fmt.Sprintf("%s/api/%s/address", lister.baseClient.Config.Host, lister.baseClient.Config.APIVersion)
		requestURL *url.URL
	)

	if requestURL, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	return requestURL, nil
}
