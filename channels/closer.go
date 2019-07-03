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

type channelCloseRequest struct {
	State string `json:"state"`
}

type Closer interface {
	Close(ctx context.Context, tokenAddress, partnerAddress common.Address) (*Channel, error)
}

func NewCloser(config *config.Config, httpClient *http.Client) Closer {
	return &defaultCloser{
		baseClient: &util.BaseClient{
			Config:     config,
			HTTPClient: httpClient,
		},
	}
}

type defaultCloser struct {
	baseClient *util.BaseClient
}

func (closer *defaultCloser) Close(ctx context.Context, tokenAddress, partnerAddress common.Address) (*Channel, error) {
	var (
		err     error
		channel = &channel{}

		requestURL          *url.URL
		request             *http.Request
		response            *http.Response
		requestBody         []byte
		channelCloseRequest = &channelCloseRequest{
			State: "closed",
		}
	)

	if requestURL, err = closer.getRequestURL(tokenAddress, partnerAddress); err != nil {
		return nil, err
	}

	if requestBody, err = json.Marshal(channelCloseRequest); err != nil {
		return nil, err
	}

	if request, err = http.NewRequest("PATCH", requestURL.String(), strings.NewReader(string(requestBody))); err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)

	if response, err = closer.baseClient.HTTPClient.Do(request); err != nil {
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

func (closer *defaultCloser) getRequestURL(tokenAddress, partnerAddress common.Address) (*url.URL, error) {
	var (
		err        error
		endpoint   = fmt.Sprintf("%s/api/%s/channels/%s/%s", closer.baseClient.Config.Host, closer.baseClient.Config.APIVersion, tokenAddress.Hex(), partnerAddress.Hex())
		requestURL *url.URL
	)

	if requestURL, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	return requestURL, nil
}
