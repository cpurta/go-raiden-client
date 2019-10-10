package tokens

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/cpurta/go-raiden-client/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExamplePartnerLister() {
	var (
		tokenClient *Client
		config      = &config.Config{
			Host:       "http://localhost:5001",
			APIVersion: "v1",
		}
		tokenAddress = common.HexToAddress("0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359") // DAI Stablecoin
		partners     []*Partner
		err          error
	)

	tokenClient = NewClient(config, http.DefaultClient)

	if partners, err = tokenClient.ListPartners(context.Background(), tokenAddress); err != nil {
		panic(fmt.Sprintf("unable to list token partners: %s", err.Error()))
	}

	fmt.Printf("token partners: %+v\n", partners)
}

func TestPartnerLister(t *testing.T) {
	var (
		localhostIP = "[::1]"
		config      = &config.Config{
			Host:       "http://localhost:5001",
			APIVersion: "v1",
		}
	)

	if os.Getenv("USE_IPV4") != "" {
		localhostIP = "127.0.0.1"
	}

	type testcase struct {
		name             string
		prepHTTPMock     func()
		expectedPartners []*Partner
		expectedError    error
	}

	testcases := []testcase{
		testcase{
			name: "successfully opened payment channel",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"GET",
					"http://localhost:5001/api/v1/tokens/0x61bB630D3B2e8eda0FC1d50F9f958eC02e3969F6/partners",
					httpmock.NewStringResponder(
						http.StatusOK,
						`[{"partner_address":"0x2a65aca4d5fc5b5c859090a6c34d164135398226","channel":"/api/v1/channels/0x61C808D82A3Ac53231750daDc13c777b59310bD9/0x2a65aca4d5fc5b5c859090a6c34d164135398226"}]`,
					),
				)
			},
			expectedError: nil,
			expectedPartners: []*Partner{
				&Partner{
					Address:    common.HexToAddress("0x2a65aca4d5fc5b5c859090a6c34d164135398226"),
					ChannelURI: "/api/v1/channels/0x61C808D82A3Ac53231750daDc13c777b59310bD9/0x2a65aca4d5fc5b5c859090a6c34d164135398226",
				},
			},
		},
		testcase{
			name: "unexpected 500 response",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"GET",
					"http://localhost:5001/api/v1/tokens/0x61bB630D3B2e8eda0FC1d50F9f958eC02e3969F6/partners",
					httpmock.NewStringResponder(
						http.StatusInternalServerError,
						``,
					),
				)
			},
			expectedError:    errors.New("EOF"),
			expectedPartners: []*Partner{},
		},
		testcase{
			name: "unable to make http request",
			prepHTTPMock: func() {
				httpmock.Deactivate()
			},
			expectedError:    fmt.Errorf("Get http://localhost:5001/api/v1/tokens/0x61bB630D3B2e8eda0FC1d50F9f958eC02e3969F6/partners: dial tcp %s:5001: connect: connection refused", localhostIP),
			expectedPartners: []*Partner{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err           error
				partners      []*Partner
				tokenAddress  = common.HexToAddress("0x61bB630D3B2e8eda0FC1d50F9f958eC02e3969F6")
				partnerLister = NewPartnerLister(config, http.DefaultClient)
				ctx           = context.Background()
			)

			httpmock.Activate()
			defer httpmock.Deactivate()

			tc.prepHTTPMock()

			partners, err = partnerLister.ListPartners(ctx, tokenAddress)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedPartners, partners)
		})
	}
}
