package payments

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/cpurta/go-raiden-client/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleLister() {
	var (
		paymentClient *Client
		config        = &config.Config{
			Host:       "http://localhost:5001",
			APIVersion: "v1",
		}
		tokenAddress  = common.HexToAddress("0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359") // DAI Stablecoin
		targetAddress = common.HexToAddress("")
		events        []*Event
		err           error
	)

	paymentClient = NewClient(config, http.DefaultClient)

	if events, err = paymentClient.List(context.Background(), tokenAddress, targetAddress); err != nil {
		panic(fmt.Sprintf("unable to list payments: %s", err.Error()))
	}

	fmt.Printf("successfully listed payment: %+v\n", events)
}

func TestLister(t *testing.T) {
	var (
		localhostIP = "[::1]"

		config = &config.Config{
			Host:       "http://localhost:5001",
			APIVersion: "v1",
		}

		time1, _ = time.Parse(time.RFC3339, "2018-10-30T07:03:52.193Z")
		time2, _ = time.Parse(time.RFC3339, "2018-10-30T07:04:22.293Z")
		time3, _ = time.Parse(time.RFC3339, "2018-10-30T07:10:13.122Z")
	)

	if os.Getenv("USE_IPV4") != "" {
		localhostIP = "127.0.0.1"
	}

	type testcase struct {
		name           string
		prepHTTPMock   func()
		expectedEvents []*Event
		expectedError  error
	}

	testcases := []testcase{
		testcase{
			name: "successfully returns at least one pending transfer",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"GET",
					"http://localhost:5001/api/v1/payments/0x0f114A1E9Db192502E7856309cc899952b3db1ED/0x82641569b2062B545431cF6D7F0A418582865ba7",
					httpmock.NewStringResponder(
						http.StatusOK,
						`[{"event":"EventPaymentReceivedSuccess","amount":5,"initiator":"0x82641569b2062B545431cF6D7F0A418582865ba7","identifier":1,"log_time":"2018-10-30T07:03:52.193Z"},{"event":"EventPaymentSentSuccess","amount":35,"target":"0x82641569b2062B545431cF6D7F0A418582865ba7","identifier":2,"log_time":"2018-10-30T07:04:22.293Z"},{"event":"EventPaymentSentSuccess","amount":20,"target":"0x82641569b2062B545431cF6D7F0A418582865ba7","identifier":3,"log_time":"2018-10-30T07:10:13.122Z"}]`,
					),
				)
			},
			expectedError: nil,
			expectedEvents: []*Event{
				&Event{
					EventName:  "EventPaymentReceivedSuccess",
					Amount:     int64(5),
					Initiator:  common.HexToAddress("0x82641569b2062B545431cF6D7F0A418582865ba7"),
					Identifier: int64(1),
					LogTime:    time1,
				},
				&Event{
					EventName:  "EventPaymentSentSuccess",
					Amount:     int64(35),
					Target:     common.HexToAddress("0x82641569b2062B545431cF6D7F0A418582865ba7"),
					Identifier: int64(2),
					LogTime:    time2,
				},
				&Event{
					EventName:  "EventPaymentSentSuccess",
					Amount:     int64(20),
					Target:     common.HexToAddress("0x82641569b2062B545431cF6D7F0A418582865ba7"),
					Identifier: int64(3),
					LogTime:    time3,
				},
			},
		},
		testcase{
			name: "unexpected 500 response",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"GET",
					"http://localhost:5001/api/v1/payments/0x0f114A1E9Db192502E7856309cc899952b3db1ED/0x82641569b2062B545431cF6D7F0A418582865ba7",
					httpmock.NewStringResponder(
						http.StatusInternalServerError,
						``,
					),
				)
			},
			expectedError:  errors.New("EOF"),
			expectedEvents: nil,
		},
		testcase{
			name: "unable to make http request",
			prepHTTPMock: func() {
				httpmock.Deactivate()
			},
			expectedError:  fmt.Errorf("Get http://localhost:5001/api/v1/payments/0x0f114A1E9Db192502E7856309cc899952b3db1ED/0x82641569b2062B545431cF6D7F0A418582865ba7: dial tcp %s:5001: connect: connection refused", localhostIP),
			expectedEvents: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err            error
				events         []*Event
				lister         = NewLister(config, http.DefaultClient)
				ctx            = context.Background()
				tokenAddress   = common.HexToAddress("0x0f114A1E9Db192502E7856309cc899952b3db1ED")
				partnerAddress = common.HexToAddress("0x82641569b2062B545431cF6D7F0A418582865ba7")
			)

			httpmock.Activate()
			defer httpmock.Deactivate()

			tc.prepHTTPMock()

			// test list all

			events, err = lister.List(ctx, tokenAddress, partnerAddress)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedEvents, events)
		})
	}
}
