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

type channelOpenRequest struct {
	PartnerAddress string `json:"partner_address"`
	TokenAddress   string `json:"token_address"`
	TotalDeposit   int64  `json:"total_deposit"`
	SettleTimeout  int64  `json:"settle_timeout"`
}

// Opener represents a generic interface to Open a Payment Channel given a token,
// partner address, deposit and a settle timeout.
type Opener interface {
	Open(ctx context.Context, tokenAddress, partnerAddress common.Address, deposit, settleTimeout int64) (*Channel, error)
}

// NewOpener creates a new default Channel opener given a Raiden node configuration
// and an http client.
func NewOpener(config *config.Config, httpClient *http.Client) Opener {
	return &defaultOpener{
		baseClient: &util.BaseClient{
			Config:     config,
			HTTPClient: httpClient,
		},
	}
}

type defaultOpener struct {
	baseClient *util.BaseClient
}

// Open will open a new payment channel given a token address, partner address, deposit, and settle timeout.
func (opener *defaultOpener) Open(ctx context.Context, tokenAddress, partnerAddress common.Address, deposit, settleTimeout int64) (*Channel, error) {
	var (
		err     error
		channel = &channel{}

		requestURL         *url.URL
		request            *http.Request
		response           *http.Response
		requestBody        []byte
		channelOpenRequest = &channelOpenRequest{
			PartnerAddress: partnerAddress.Hex(),
			TokenAddress:   tokenAddress.Hex(),
			TotalDeposit:   deposit,
			SettleTimeout:  settleTimeout,
		}
	)

	if requestURL, err = opener.getRequestURL(); err != nil {
		return nil, err
	}

	if requestBody, err = json.Marshal(channelOpenRequest); err != nil {
		return nil, err
	}

	if request, err = http.NewRequest("PUT", requestURL.String(), strings.NewReader(string(requestBody))); err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)

	if response, err = opener.baseClient.HTTPClient.Do(request); err != nil {
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

func (opener *defaultOpener) getRequestURL() (*url.URL, error) {
	var (
		err        error
		endpoint   = fmt.Sprintf("%s/api/%s/channels", opener.baseClient.Config.Host, opener.baseClient.Config.APIVersion)
		requestURL *url.URL
	)

	if requestURL, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	return requestURL, nil
}
