package payments

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

func TestInitiator(t *testing.T) {
	var (
		config = &config.Config{
			Host:       "http://localhost:5001",
			APIVersion: "v1",
		}
	)

	type testcase struct {
		name            string
		prepHTTPMock    func()
		expectedPayment *Payment
		expectedError   error
	}

	testcases := []testcase{
		testcase{
			name: "successfully initiated a funds transfer",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"POST",
					"http://localhost:5001/api/v1/payments/0x2a65Aca4D5fC5B5C859090a6c34d164135398226/0x61C808D82A3Ac53231750daDc13c777b59310bD9",
					httpmock.NewStringResponder(
						http.StatusOK,
						`{"initiator_address":"0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8","target_address":"0x61C808D82A3Ac53231750daDc13c777b59310bD9","token_address":"0x2a65Aca4D5fC5B5C859090a6c34d164135398226","amount":200,"identifier":42}`,
					),
				)
			},
			expectedError: nil,
			expectedPayment: &Payment{
				InitiatorAddress: common.HexToAddress("0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8"),
				TargetAddress:    common.HexToAddress("0x61C808D82A3Ac53231750daDc13c777b59310bD9"),
				TokenAddress:     common.HexToAddress("0x2a65Aca4D5fC5B5C859090a6c34d164135398226"),
				Amount:           int64(200),
				Identifier:       int64(42),
			},
		},
		testcase{
			name: "unexpected 500 response",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"POST",
					"http://localhost:5001/api/v1/payments/0x2a65Aca4D5fC5B5C859090a6c34d164135398226/0x61C808D82A3Ac53231750daDc13c777b59310bD9",
					httpmock.NewStringResponder(
						http.StatusInternalServerError,
						``,
					),
				)
			},
			expectedError:   errors.New("EOF"),
			expectedPayment: nil,
		},
		testcase{
			name: "unable to make http request",
			prepHTTPMock: func() {
				httpmock.Deactivate()
			},
			expectedError:   errors.New("Post http://localhost:5001/api/v1/payments/0x2a65Aca4D5fC5B5C859090a6c34d164135398226/0x61C808D82A3Ac53231750daDc13c777b59310bD9: dial tcp [::1]:5001: connect: connection refused"),
			expectedPayment: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err            error
				payment        *Payment
				initiator      = NewInitiator(config, http.DefaultClient)
				ctx            = context.Background()
				tokenAddress   = common.HexToAddress("0x2a65Aca4D5fC5B5C859090a6c34d164135398226")
				partnerAddress = common.HexToAddress("0x61C808D82A3Ac53231750daDc13c777b59310bD9")
			)

			httpmock.Activate()
			defer httpmock.Deactivate()

			tc.prepHTTPMock()

			// test list all

			payment, err = initiator.Initiate(ctx, tokenAddress, partnerAddress, int64(200))

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedPayment, payment)
		})
	}
}
