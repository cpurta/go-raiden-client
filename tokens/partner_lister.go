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

// PartnerLister is a generic interface to list the Ethereum address associated to the
// Raiden node. It allows for a context to be passed to allow for request timeouts
// and/or deadlines on the response.
type PartnerLister interface {
	ListPartners(ctx context.Context, tokenAddress common.Address) ([]*Partner, error)
}

var _ PartnerLister = &defaultPartnerLister{}

// NewPartnerLister will return a default address lister for a configured Raiden node.
func NewPartnerLister(config *config.Config, httpClient *http.Client) PartnerLister {
	return &defaultPartnerLister{
		baseClient: &util.BaseClient{
			Config:     config,
			HTTPClient: httpClient,
		},
	}
}

type defaultPartnerLister struct {
	baseClient *util.BaseClient
}

// ListPartners will return the associated Ethereum address to the Raiden node configured
// in the Lister config.
func (lister *defaultPartnerLister) ListPartners(ctx context.Context, tokenAddress common.Address) ([]*Partner, error) {
	var (
		err      error
		partners = make([]*Partner, 0)

		requestURL *url.URL
		request    *http.Request
		response   *http.Response
	)

	if requestURL, err = lister.getRequestURL(tokenAddress); err != nil {
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

	if err = json.NewDecoder(response.Body).Decode(&partners); err != nil {
		return nil, err
	}

	return partners, nil
}

func (lister *defaultPartnerLister) getRequestURL(tokenAddress common.Address) (*url.URL, error) {
	var (
		err        error
		endpoint   = fmt.Sprintf("%s/api/%s/tokens/%s/partners", lister.baseClient.Config.Host, lister.baseClient.Config.APIVersion, tokenAddress.Hex())
		requestURL *url.URL
	)

	if requestURL, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	return requestURL, nil
}
