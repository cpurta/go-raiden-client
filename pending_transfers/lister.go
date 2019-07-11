package pendingtransfers

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

var (
	zeroAddress        = common.Address{}
	_           Lister = &defaultLister{}
)

// Lister is an interface that allows for various list operations to be performed.
type Lister interface {
	ListAll(context.Context) ([]*Transfer, error)
	ListToken(context.Context, common.Address) ([]*Transfer, error)
	ListChannel(context.Context, common.Address, common.Address) ([]*Transfer, error)
}

// NewLister will return a default lister that will be able to perform the various
// listing operations of all pending transfers know by a Raiden node.
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

// ListAll will list all currently pending transfers on the Raiden node.
func (lister *defaultLister) ListAll(ctx context.Context) ([]*Transfer, error) {
	var (
		url *url.URL
		err error
	)

	if url, err = lister.getAllRequestURL(); err != nil {
		return nil, err
	}

	return lister.getPendingTransfers(ctx, url)
}

func (lister *defaultLister) ListToken(ctx context.Context, tokenAddress common.Address) ([]*Transfer, error) {
	var (
		url *url.URL
		err error
	)

	if url, err = lister.getTokenRequestURL(tokenAddress); err != nil {
		return nil, err
	}

	return lister.getPendingTransfers(ctx, url)
}

func (lister *defaultLister) ListChannel(ctx context.Context, tokenAddress common.Address, partnerAddress common.Address) ([]*Transfer, error) {
	var (
		url *url.URL
		err error
	)

	if url, err = lister.getChannelRequestURL(tokenAddress, partnerAddress); err != nil {
		return nil, err
	}

	return lister.getPendingTransfers(ctx, url)
}

func (lister *defaultLister) getPendingTransfers(ctx context.Context, url *url.URL) ([]*Transfer, error) {
	var (
		err       error
		transfers = make([]*Transfer, 0)

		request  *http.Request
		response *http.Response
	)

	if request, err = http.NewRequest("GET", url.String(), nil); err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)

	if response, err = lister.baseClient.HTTPClient.Do(request); err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&transfers); err != nil {
		return nil, err
	}

	return transfers, nil
}

func (lister *defaultLister) getAllRequestURL() (*url.URL, error) {
	var (
		err        error
		endpoint   = fmt.Sprintf("%s/api/%s/pending_transfers", lister.baseClient.Config.Host, lister.baseClient.Config.APIVersion)
		requestURL *url.URL
	)

	if requestURL, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	return requestURL, nil
}

func (lister *defaultLister) getChannelRequestURL(tokenAddress, partnerAddress common.Address) (*url.URL, error) {
	var (
		err        error
		endpoint   = fmt.Sprintf("%s/api/%s/pending_transfers/%s/%s", lister.baseClient.Config.Host, lister.baseClient.Config.APIVersion, tokenAddress.Hex(), partnerAddress.Hex())
		requestURL *url.URL
	)

	if requestURL, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	return requestURL, nil
}

func (lister *defaultLister) getTokenRequestURL(tokenAddress common.Address) (*url.URL, error) {
	var (
		err        error
		endpoint   = fmt.Sprintf("%s/api/%s/pending_transfers/%s", lister.baseClient.Config.Host, lister.baseClient.Config.APIVersion, tokenAddress.Hex())
		requestURL *url.URL
	)

	if requestURL, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	return requestURL, nil
}
