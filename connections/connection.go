package connections

// Connection represents a high level information about the Funds, Total deposits
// and numbers of channels for a given network.
type Connection struct {
	Funds       int64 `json:"funds"`
	SumDeposits int64 `json:"sum_deposits"`
	Channels    int64 `json:"channels"`
}
