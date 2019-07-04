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

func TestLister(t *testing.T) {
	var (
		config = &config.Config{
			Host:       "http://localhost:5001",
			APIVersion: "v1",
		}
	)

	type testcase struct {
		name                string
		prepHTTPMock        func()
		expectedConnections Connections
		expectedError       error
	}

	testcases := []testcase{
		testcase{
			name: "successfully joined a token network",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"GET",
					"http://localhost:5001/api/v1/connections",
					httpmock.NewStringResponder(
						http.StatusNoContent,
						`{"0x2a65Aca4D5fC5B5C859090a6c34d164135398226":{"funds":100,"sum_deposits":67,"channels":3},"0x0f114A1E9Db192502E7856309cc899952b3db1ED":{"funds":49,"sum_deposits":31,"channels":1}}`,
					),
				)
			},
			expectedConnections: Connections{
				common.HexToAddress("0x2a65Aca4D5fC5B5C859090a6c34d164135398226"): &Connection{
					Funds:       int64(100),
					SumDeposits: int64(67),
					Channels:    int64(3),
				},
				common.HexToAddress("0x0f114A1E9Db192502E7856309cc899952b3db1ED"): &Connection{
					Funds:       int64(49),
					SumDeposits: int64(31),
					Channels:    int64(1),
				},
			},
			expectedError: nil,
		},
		testcase{
			name: "unexpected 500 response",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"GET",
					"http://localhost:5001/api/v1/connections",
					httpmock.NewStringResponder(
						http.StatusInternalServerError,
						``,
					),
				)
			},
			expectedConnections: nil,
			expectedError:       errors.New("EOF"),
		},
		testcase{
			name: "unable to make http request",
			prepHTTPMock: func() {
				httpmock.Deactivate()
			},
			expectedConnections: nil,
			expectedError:       errors.New("Get http://localhost:5001/api/v1/connections: dial tcp [::1]:5001: connect: connection refused"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err         error
				connections Connections

				lister = NewLister(config, http.DefaultClient)
				ctx    = context.Background()
			)

			httpmock.Activate()
			defer httpmock.Deactivate()

			tc.prepHTTPMock()

			connections, err = lister.List(ctx)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedConnections, connections)
		})
	}
}
