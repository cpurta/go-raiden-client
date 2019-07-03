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

// Lister is a generic interface to list the Ethereum address associated to the
// Raiden node. It allows for a context to be passed to allow for request timeouts
// and/or deadlines on the response.
type Lister interface {
	List(ctx context.Context) ([]common.Address, error)
}

var _ Lister = &defaultLister{}

// NewLister will return a default address lister for a configured Raiden node.
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

// List will return the associated Ethereum address to the Raiden node configured
// in the Lister config.
func (lister *defaultLister) List(ctx context.Context) ([]common.Address, error) {
	var (
		err       error
		addresses = make([]common.Address, 0)

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

	if err = json.NewDecoder(response.Body).Decode(&addresses); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (lister *defaultLister) getRequestURL() (*url.URL, error) {
	var (
		err        error
		endpoint   = fmt.Sprintf("%s/api/%s/tokens", lister.baseClient.Config.Host, lister.baseClient.Config.APIVersion)
		requestURL *url.URL
	)

	if requestURL, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	return requestURL, nil
}
