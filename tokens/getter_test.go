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

func ExampleGetter() {
	var (
		tokenClient *Client
		config      = &config.Config{
			Host:       "http://localhost:5001",
			APIVersion: "v1",
		}
		tokenAddress = common.HexToAddress("0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359") // DAI Stablecoin
		address      common.Address
		err          error
	)

	tokenClient = NewClient(config, http.DefaultClient)

	if address, err = tokenClient.Get(context.Background(), tokenAddress); err != nil {
		panic(fmt.Sprintf("unable to get address: %s", err.Error()))
	}

	fmt.Printf("address: %s\n", address.String())
}

func TestGetter(t *testing.T) {
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
		name            string
		prepHTTPMock    func()
		expectedAddress common.Address
		expectedError   error
	}

	testcases := []testcase{
		testcase{
			name: "successfully opened payment channel",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"GET",
					"http://localhost:5001/api/v1/tokens/0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8",
					httpmock.NewStringResponder(
						http.StatusOK,
						`"0x61bB630D3B2e8eda0FC1d50F9f958eC02e3969F6"`,
					),
				)
			},
			expectedError:   nil,
			expectedAddress: common.HexToAddress("0x61bB630D3B2e8eda0FC1d50F9f958eC02e3969F6"),
		},
		testcase{
			name: "unexpected 500 response",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"GET",
					"http://localhost:5001/api/v1/tokens/0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8",
					httpmock.NewStringResponder(
						http.StatusInternalServerError,
						``,
					),
				)
			},
			expectedError:   errors.New("EOF"),
			expectedAddress: common.Address{},
		},
		testcase{
			name: "unable to make http request",
			prepHTTPMock: func() {
				httpmock.Deactivate()
			},
			expectedError:   fmt.Errorf("Get http://localhost:5001/api/v1/tokens/0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8: dial tcp %s:5001: connect: connection refused", localhostIP),
			expectedAddress: common.Address{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err          error
				address      common.Address
				tokenAddress = common.HexToAddress("0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8")

				getter = NewGetter(config, http.DefaultClient)
				ctx    = context.Background()
			)

			httpmock.Activate()
			defer httpmock.Deactivate()

			tc.prepHTTPMock()

			address, err = getter.Get(ctx, tokenAddress)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedAddress, address)
		})
	}
}
