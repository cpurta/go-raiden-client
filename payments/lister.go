package payments

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/cpurta/go-raiden-client/config"
	"github.com/cpurta/go-raiden-client/util"
	"github.com/ethereum/go-ethereum/common"
)

type Lister interface {
	List(ctx context.Context, tokenAddress, targetAddress common.Address) ([]*Event, error)
}

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

func (lister *defaultLister) List(ctx context.Context, tokenAddress, targetAddress common.Address) ([]*Event, error) {
	var (
		err           error
		events        = make([]*event, 0)
		paymentEvents = make([]*Event, 0)

		requestURL *url.URL
		request    *http.Request
		response   *http.Response
	)

	if requestURL, err = lister.getRequestURL(tokenAddress, targetAddress); err != nil {
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

	if err = json.NewDecoder(response.Body).Decode(&events); err != nil {
		return nil, err
	}

	for _, event := range events {
		var (
			logTime time.Time
		)
		if logTime, err = time.Parse(time.RFC3339, event.LogTime); err != nil {
			continue
		}

		paymentEvents = append(paymentEvents, &Event{
			EventName:  event.EventName,
			Amount:     event.Amount,
			Initiator:  common.HexToAddress(event.Initiator),
			Identifier: event.Identifier,
			LogTime:    logTime,
		})
	}

	return paymentEvents, nil
}

func (lister *defaultLister) getRequestURL(tokenAddress, targetAddress common.Address) (*url.URL, error) {
	var (
		err        error
		endpoint   = fmt.Sprintf("%s/api/%s/payments/%s/%s", lister.baseClient.Config.Host, lister.baseClient.Config.APIVersion, tokenAddress.Hex(), targetAddress.Hex())
		requestURL *url.URL
	)

	if requestURL, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	return requestURL, nil
}
