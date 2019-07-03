package util

import (
	"net/http"

	"github.com/cpurta/go-raiden-client/config"
)

// BaseClient serves as the HTTP client responsible for making all outbound requests
// to the Raiden node as specified in the Config. It allows for HTTP requests to be
// built onto and add headers if needed.
type BaseClient struct {
	Config     *config.Config
	HTTPClient *http.Client
}
