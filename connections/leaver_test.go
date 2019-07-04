package connections

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

func TestLeaver(t *testing.T) {
	var (
		config = &config.Config{
			Host:       "http://localhost:5001",
			APIVersion: "v1",
		}
	)

	type testcase struct {
		name              string
		prepHTTPMock      func()
		expectedAddresses []common.Address
		expectedError     error
	}

	testcases := []testcase{
		testcase{
			name: "successfully joined a token network",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"DELETE",
					"http://localhost:5001/api/v1/connections/0x2a65Aca4D5fC5B5C859090a6c34d164135398226",
					httpmock.NewStringResponder(
						http.StatusNoContent,
						`["0x41BCBC2fD72a731bcc136Cf6F7442e9C19e9f313","0x5A5f458F6c1a034930E45dC9a64B99d7def06D7E","0x8942c06FaA74cEBFf7d55B79F9989AdfC85C6b85"]`,
					),
				)
			},
			expectedAddresses: []common.Address{
				common.HexToAddress("0x41BCBC2fD72a731bcc136Cf6F7442e9C19e9f313"),
				common.HexToAddress("0x5A5f458F6c1a034930E45dC9a64B99d7def06D7E"),
				common.HexToAddress("0x8942c06FaA74cEBFf7d55B79F9989AdfC85C6b85"),
			},
			expectedError: nil,
		},
		testcase{
			name: "unexpected 500 response",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"DELETE",
					"http://localhost:5001/api/v1/connections/0x2a65Aca4D5fC5B5C859090a6c34d164135398226",
					httpmock.NewStringResponder(
						http.StatusInternalServerError,
						``,
					),
				)
			},
			expectedAddresses: []common.Address{},
			expectedError:     errors.New("EOF"),
		},
		testcase{
			name: "unable to make http request",
			prepHTTPMock: func() {
				httpmock.Deactivate()
			},
			expectedAddresses: []common.Address{},
			expectedError:     errors.New("Delete http://localhost:5001/api/v1/connections/0x2a65Aca4D5fC5B5C859090a6c34d164135398226: dial tcp [::1]:5001: connect: connection refused"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err       error
				addresses []common.Address

				tokenAddress = common.HexToAddress("0x2a65Aca4D5fC5B5C859090a6c34d164135398226")
				leaver       = NewLeaver(config, http.DefaultClient)
				ctx          = context.Background()
			)

			httpmock.Activate()
			defer httpmock.Deactivate()

			tc.prepHTTPMock()

			addresses, err = leaver.Leave(ctx, tokenAddress)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedAddresses, addresses)
		})
	}
}
