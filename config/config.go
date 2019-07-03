package config

// Config holds the needed information for a Raiden client to make API requests
// to a Raiden node.
type Config struct {
	Host       string
	APIVersion string
}
