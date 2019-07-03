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

type Leaver interface {
	Leave(ctx context.Context, tokenAddress common.Address) ([]common.Address, error)
}

func NewLeaver(config *config.Config, httpClient *http.Client) Leaver {
	return &defaultLeaver{
		baseClient: &util.BaseClient{
			Config:     config,
			HTTPClient: httpClient,
		},
	}
}

type defaultLeaver struct {
	baseClient *util.BaseClient
}

func (leaver *defaultLeaver) Leave(ctx context.Context, tokenAddress common.Address) ([]common.Address, error) {
	var (
		err            error
		tokens         = make([]string, 0)
		tokenAddresses = make([]common.Address, 0)

		requestURL *url.URL
		request    *http.Request
		response   *http.Response
	)

	if requestURL, err = leaver.getRequestURL(tokenAddress); err != nil {
		return nil, err
	}

	if request, err = http.NewRequest("DELETE", requestURL.String(), nil); err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)

	if response, err = leaver.baseClient.HTTPClient.Do(request); err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&tokens); err != nil {
		return nil, err
	}

	for _, token := range tokens {
		tokenAddresses = append(tokenAddresses, common.HexToAddress(token))
	}

	return tokenAddresses, nil
}

func (leaver *defaultLeaver) getRequestURL(tokenAddress common.Address) (*url.URL, error) {
	var (
		err        error
		endpoint   = fmt.Sprintf("%s/api/%s/connections/%s", leaver.baseClient.Config.Host, leaver.baseClient.Config.APIVersion, tokenAddress.Hex())
		requestURL *url.URL
	)

	if requestURL, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	return requestURL, nil
}
