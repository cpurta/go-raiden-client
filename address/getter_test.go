package address

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/cpurta/go-raiden-client/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGetter(t *testing.T) {
	var (
		config = &config.Config{
			Host:       "http://localhost:5001",
			APIVersion: "v1",
		}
	)

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
			expectedError:   errors.New("Get http://localhost:5001/api/v1/address: dial tcp [::1]:5001: connect: connection refused"),
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
