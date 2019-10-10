package address

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
		addressClient *Client
		config        = &config.Config{
			Host:       "http://localhost:5001",
			APIVersion: "v1",
		}
		address common.Address
		err     error
	)

	addressClient = NewClient(config, http.DefaultClient)

	if address, err = addressClient.Get(context.Background()); err != nil {
		panic(fmt.Sprintf("unable to get ethereum address from raiden node: %s", err.Error()))
	}

	fmt.Println("raiden address:", address.String())
}

func TestNewGetter(t *testing.T) {
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
			name: "successfully got ethereum address",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"GET",
					"http://localhost:5001/api/v1/address",
					httpmock.NewStringResponder(
						http.StatusOK,
						`{"our_address":"0x2a65Aca4D5fC5B5C859090a6c34d164135398226"}`,
					),
				)
			},
			expectedError:   nil,
			expectedAddress: common.HexToAddress("0x2a65Aca4D5fC5B5C859090a6c34d164135398226"),
		},
		testcase{
			name: "unexpected 500 response",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"GET",
					"http://localhost:5001/api/v1/address",
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
			expectedError:   fmt.Errorf("Get http://localhost:5001/api/v1/address: dial tcp %s:5001: connect: connection refused", localhostIP),
			expectedAddress: common.Address{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err     error
				address common.Address

				getter = NewGetter(config, http.DefaultClient)
				ctx    = context.Background()
			)

			httpmock.Activate()
			defer httpmock.Deactivate()

			tc.prepHTTPMock()

			address, err = getter.Get(ctx)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedAddress, address)
		})
	}
}
