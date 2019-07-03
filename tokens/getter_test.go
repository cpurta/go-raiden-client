package tokens

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

func TestGetter(t *testing.T) {
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
			expectedError:   errors.New("Get http://localhost:5001/api/v1/tokens/0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8: dial tcp [::1]:5001: connect: connection refused"),
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