package pendingtransfers

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

func ExampleLister() {
	var (
		transfersClient *Client
		config          = &config.Config{
			Host:       "http://localhost:5001",
			APIVersion: "v1",
		}
		tokenAddress   = common.HexToAddress("0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359") // DAI Stablecoin
		partnerAddress = common.HexToAddress("")
		transfers      []*Transfer
		err            error
	)

	transfersClient = NewClient(config, http.DefaultClient)

	if transfers, err = transfersClient.ListAll(context.Background()); err != nil {
		panic(fmt.Sprintf("unable to list all pending transfers: %s", err.Error()))
	}

	fmt.Printf("all pending transfers: %+v\n", transfers)

	if transfers, err = transfersClient.ListToken(context.Background(), tokenAddress); err != nil {
		panic(fmt.Sprintf("unable to list token pending transfers: %s", err.Error()))
	}

	fmt.Printf("token pending transfers: %+v\n", transfers)

	if transfers, err = transfersClient.ListChannel(context.Background(), tokenAddress, partnerAddress); err != nil {
		panic(fmt.Sprintf("unable to list channel pending transfers: %s", err.Error()))
	}

	fmt.Printf("channel pending transfers: %+v\n", transfers)
}

func TestLister(t *testing.T) {
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
		name              string
		prepHTTPMock      func()
		expectedTransfers []*Transfer
		expectedError     error
	}

	testcases := []testcase{
		testcase{
			name: "successfully returns at least one pending transfer",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"GET",
					"http://localhost:5001/api/v1/pending_transfers",
					httpmock.NewStringResponder(
						http.StatusOK,
						`[{"channel_identifier":255,"initiator":"0x5E1a3601538f94c9e6D2B40F7589030ac5885FE7","locked_amount":119,"payment_identifier":1,"role":"initiator","target":"0x00AF5cBfc8dC76cd599aF623E60F763228906F3E","token_address":"0xd0A1E359811322d97991E03f863a0C30C2cF029C","token_network_identifier":"0x111157460c0F41EfD9107239B7864c062aA8B978","transferred_amount":331}]`,
					),
				)

				httpmock.RegisterResponder(
					"GET",
					"http://localhost:5001/api/v1/pending_transfers/0xd0A1E359811322d97991E03f863a0C30C2cF029C",
					httpmock.NewStringResponder(
						http.StatusOK,
						`[{"channel_identifier":255,"initiator":"0x5E1a3601538f94c9e6D2B40F7589030ac5885FE7","locked_amount":119,"payment_identifier":1,"role":"initiator","target":"0x00AF5cBfc8dC76cd599aF623E60F763228906F3E","token_address":"0xd0A1E359811322d97991E03f863a0C30C2cF029C","token_network_identifier":"0x111157460c0F41EfD9107239B7864c062aA8B978","transferred_amount":331}]`,
					),
				)

				httpmock.RegisterResponder(
					"GET",
					"http://localhost:5001/api/v1/pending_transfers/0xd0A1E359811322d97991E03f863a0C30C2cF029C/0x2c4b0Bdac486d492E3cD701F4cA87e480AE4C685",
					httpmock.NewStringResponder(
						http.StatusOK,
						`[{"channel_identifier":255,"initiator":"0x5E1a3601538f94c9e6D2B40F7589030ac5885FE7","locked_amount":119,"payment_identifier":1,"role":"initiator","target":"0x00AF5cBfc8dC76cd599aF623E60F763228906F3E","token_address":"0xd0A1E359811322d97991E03f863a0C30C2cF029C","token_network_identifier":"0x111157460c0F41EfD9107239B7864c062aA8B978","transferred_amount":331}]`,
					),
				)
			},
			expectedError: nil,
			expectedTransfers: []*Transfer{
				&Transfer{
					ChannelIdentifier:      int64(255),
					Initiator:              common.HexToAddress("0x5E1a3601538f94c9e6D2B40F7589030ac5885FE7"),
					LockedAmount:           int64(119),
					PaymentIdentifier:      int64(1),
					Role:                   "initiator",
					Target:                 common.HexToAddress("0x00AF5cBfc8dC76cd599aF623E60F763228906F3E"),
					TokenAddress:           common.HexToAddress("0xd0A1E359811322d97991E03f863a0C30C2cF029C"),
					TokenNetworkIdentifier: common.HexToAddress("0x111157460c0F41EfD9107239B7864c062aA8B978"),
					TransferredAmount:      int64(331),
				},
			},
		},
		testcase{
			name: "unexpected 500 response",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"GET",
					"http://localhost:5001/api/v1/pending_transfers",
					httpmock.NewStringResponder(
						http.StatusInternalServerError,
						``,
					),
				)

				httpmock.RegisterResponder(
					"GET",
					"http://localhost:5001/api/v1/pending_transfers/0xd0A1E359811322d97991E03f863a0C30C2cF029C",
					httpmock.NewStringResponder(
						http.StatusInternalServerError,
						``,
					),
				)

				httpmock.RegisterResponder(
					"GET",
					"http://localhost:5001/api/v1/pending_transfers/0xd0A1E359811322d97991E03f863a0C30C2cF029C/0x2c4b0Bdac486d492E3cD701F4cA87e480AE4C685",
					httpmock.NewStringResponder(
						http.StatusInternalServerError,
						``,
					),
				)
			},
			expectedError:     errors.New("EOF"),
			expectedTransfers: nil,
		},
		testcase{
			name: "unable to make http request",
			prepHTTPMock: func() {
				httpmock.Deactivate()
			},
			expectedError:     fmt.Errorf("Get http://localhost:5001/api/v1/pending_transfers: dial tcp %s:5001: connect: connection refused", localhostIP),
			expectedTransfers: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err            error
				transfers      []*Transfer
				lister         = NewLister(config, http.DefaultClient)
				ctx            = context.Background()
				tokenAddress   = common.HexToAddress("0xd0A1E359811322d97991E03f863a0C30C2cF029C")
				partnerAddress = common.HexToAddress("0x2c4b0Bdac486d492E3cD701F4cA87e480AE4C685")
			)

			httpmock.Activate()
			defer httpmock.Deactivate()

			tc.prepHTTPMock()

			// test list all

			transfers, err = lister.ListAll(ctx)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedTransfers, transfers)

			// test token filtered

			transfers, err = lister.ListToken(ctx, tokenAddress)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			assert.Equal(t, tc.expectedTransfers, transfers)

			// test channel filtered

			transfers, err = lister.ListChannel(ctx, tokenAddress, partnerAddress)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			assert.Equal(t, tc.expectedTransfers, transfers)
		})
	}
}
