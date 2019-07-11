package connections

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/cpurta/go-raiden-client/config"
	"github.com/cpurta/go-raiden-client/util"
	"github.com/ethereum/go-ethereum/common"
)

type joinRequest struct {
	Funds int64 `json:"funds"`
}

// Joiner is an interface to allow for a Raiden node to join a new token network
// with a given number of funds.
type Joiner interface {
	Join(ctx context.Context, tokenAddress common.Address, funds int64) error
}

// NewJoiner will create a default joiner that will allow access to join a new token
// network for a Raiden node.
func NewJoiner(config *config.Config, httpClient *http.Client) Joiner {
	return &defaultJoiner{
		baseClient: &util.BaseClient{
			Config:     config,
			HTTPClient: httpClient,
		},
	}
}

type defaultJoiner struct {
	baseClient *util.BaseClient
}

// Join will join a new token network given a token network address and a given number of funds.
func (joiner *defaultJoiner) Join(ctx context.Context, tokenAddress common.Address, funds int64) error {
	var (
		err          error
		requestURL   *url.URL
		request      *http.Request
		response     *http.Response
		requestBody  []byte
		responseBody []byte
		joinRequest  = &joinRequest{
			Funds: funds,
		}
	)

	if requestURL, err = joiner.getRequestURL(tokenAddress); err != nil {
		return err
	}

	if requestBody, err = json.Marshal(joinRequest); err != nil {
		return err
	}

	if request, err = http.NewRequest("PUT", requestURL.String(), strings.NewReader(string(requestBody))); err != nil {
		return err
	}

	request = request.WithContext(ctx)

	if response, err = joiner.baseClient.HTTPClient.Do(request); err != nil {
		return err
	}

	defer response.Body.Close()

	if responseBody, err = ioutil.ReadAll(response.Body); err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("recieved %d status code: %s", response.StatusCode, string(responseBody))
	}

	return nil
}

func (joiner *defaultJoiner) getRequestURL(tokenAddress common.Address) (*url.URL, error) {
	var (
		err        error
		endpoint   = fmt.Sprintf("%s/api/%s/connections/%s", joiner.baseClient.Config.Host, joiner.baseClient.Config.APIVersion, tokenAddress.Hex())
		requestURL *url.URL
	)

	if requestURL, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	return requestURL, nil
}
