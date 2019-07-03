package payments

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

type initiatePaymentRequest struct {
	Amount int64 `json:"amount"`
}

type Initiator interface {
	Initiate(ctx context.Context, tokenAddress, targetAddress common.Address, amount int64) (*Payment, error)
}

func NewInitiator(config *config.Config, httpClient *http.Client) Initiator {
	return &defaultInitiator{
		baseClient: &util.BaseClient{
			Config:     config,
			HTTPClient: httpClient,
		},
	}
}

type defaultInitiator struct {
	baseClient *util.BaseClient
}

func (initiator *defaultInitiator) Initiate(ctx context.Context, tokenAddress, targetAddress common.Address, amount int64) (*Payment, error) {
	var (
		err     error
		payment *Payment

		requestURL *url.URL
		request    *http.Request
		response   *http.Response
	)

	if requestURL, err = initiator.getRequestURL(tokenAddress, targetAddress); err != nil {
		return nil, err
	}

	if request, err = http.NewRequest("POST", requestURL.String(), nil); err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)

	if response, err = initiator.baseClient.HTTPClient.Do(request); err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (initiator *defaultInitiator) getRequestURL(tokenAddress, targetAddress common.Address) (*url.URL, error) {
	var (
		err        error
		endpoint   = fmt.Sprintf("%s/api/%s/payments/%s/%s", initiator.baseClient.Config.Host, initiator.baseClient.Config.APIVersion, tokenAddress.Hex(), targetAddress.Hex())
		requestURL *url.URL
	)

	if requestURL, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	return requestURL, nil
}
