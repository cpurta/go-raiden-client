package channels

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/cpurta/go-raiden-client/config"
	"github.com/cpurta/go-raiden-client/util"
	"github.com/ethereum/go-ethereum/common"
)

type increaseDepositRequest struct {
	TotalDeposit int64 `json:"total_deposit"`
}

// IncreaseDepositor represents a generic interface to Increase the Deposit of a Payment Channel given a token and
// partner address.
type IncreaseDepositor interface {
	IncreaseDeposit(ctx context.Context, tokenAddress, partnerAddress common.Address, deposit int64) (*Channel, error)
}

// NewIncreaseDepositor creates a new default Channel depositor increaser given a Raiden node configuration
// and an http client.
func NewIncreaseDepositor(config *config.Config, httpClient *http.Client) IncreaseDepositor {
	return &defaultIncreaseDepositor{
		baseClient: &util.BaseClient{
			Config:     config,
			HTTPClient: httpClient,
		},
	}
}

type defaultIncreaseDepositor struct {
	baseClient *util.BaseClient
}

// Close will increase the deposit a payment channel given a token address and a partner address.
func (depositor *defaultIncreaseDepositor) IncreaseDeposit(ctx context.Context, tokenAddress, partnerAddress common.Address, deposit int64) (*Channel, error) {
	var (
		err     error
		channel = &channel{}

		requestURL             *url.URL
		request                *http.Request
		response               *http.Response
		requestBody            []byte
		increaseDepositRequest = &increaseDepositRequest{
			TotalDeposit: deposit,
		}
	)

	if requestURL, err = depositor.getRequestURL(tokenAddress, partnerAddress); err != nil {
		return nil, err
	}

	if requestBody, err = json.Marshal(increaseDepositRequest); err != nil {
		return nil, err
	}

	if request, err = http.NewRequest("PATCH", requestURL.String(), strings.NewReader(string(requestBody))); err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)

	if response, err = depositor.baseClient.HTTPClient.Do(request); err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&channel); err != nil {
		return nil, err
	}

	return &Channel{
		TokenNetworkIdentifier: common.HexToAddress(channel.TokenNetworkIdentifier),
		ChannelIdentifier:      channel.ChannelIdentifier,
		PartnerAddress:         common.HexToAddress(channel.PartnerAddress),
		TokenAddress:           common.HexToAddress(channel.TokenAddress),
		Balance:                channel.Balance,
		TotalDeposit:           channel.TotalDeposit,
		State:                  channel.State,
		SettleTimeout:          channel.SettleTimeout,
		RevealTimeout:          channel.RevealTimeout,
	}, nil
}

func (depositor *defaultIncreaseDepositor) getRequestURL(tokenAddress, partnerAddress common.Address) (*url.URL, error) {
	var (
		err        error
		endpoint   = fmt.Sprintf("%s/api/%s/channels/%s/%s", depositor.baseClient.Config.Host, depositor.baseClient.Config.APIVersion, tokenAddress.Hex(), partnerAddress.Hex())
		requestURL *url.URL
	)

	if requestURL, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	return requestURL, nil
}
