package connections

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

func ExampleJoiner() {
	var (
		connClient *Client
		config     = &config.Config{
			Host:       "http://localhost:5001",
			APIVersion: "v1",
		}
		tokenAddress = common.HexToAddress("0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359") // DAI Stablecoin
		funds        = int64(1000)
		err          error
	)

	connClient = NewClient(config, http.DefaultClient)

	if err = connClient.Join(context.Background(), tokenAddress, funds); err != nil {
		panic(fmt.Sprintf("unable to join connection: %s", err.Error()))
	}

	fmt.Printf("succuessfully join connection")
}

func TestCloser(t *testing.T) {
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
		name          string
		prepHTTPMock  func()
		expectedError error
	}

	testcases := []testcase{
		testcase{
			name: "successfully joined a token network",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"PUT",
					"http://localhost:5001/api/v1/connections/0x2a65Aca4D5fC5B5C859090a6c34d164135398226",
					httpmock.NewStringResponder(
						http.StatusNoContent,
						`{"funds":1337}`,
					),
				)
			},
			expectedError: nil,
		},
		testcase{
			name: "unexpected 500 response",
			prepHTTPMock: func() {
				httpmock.RegisterResponder(
					"PUT",
					"http://localhost:5001/api/v1/connections/0x2a65Aca4D5fC5B5C859090a6c34d164135398226",
					httpmock.NewStringResponder(
						http.StatusInternalServerError,
						``,
					),
				)
			},
			expectedError: errors.New("recieved 500 status code: "),
		},
		testcase{
			name: "unable to make http request",
			prepHTTPMock: func() {
				httpmock.Deactivate()
			},
			expectedError: fmt.Errorf("Put http://localhost:5001/api/v1/connections/0x2a65Aca4D5fC5B5C859090a6c34d164135398226: dial tcp %s:5001: connect: connection refused", localhostIP),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err          error
				tokenAddress = common.HexToAddress("0x2a65Aca4D5fC5B5C859090a6c34d164135398226")

				joiner = NewJoiner(config, http.DefaultClient)
				ctx    = context.Background()
			)

			httpmock.Activate()
			defer httpmock.Deactivate()

			tc.prepHTTPMock()

			err = joiner.Join(ctx, tokenAddress, 1337)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}
