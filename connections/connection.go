package connections

type Connection struct {
	Funds       int64 `json:"funds"`
	SumDeposits int64 `json:"sum_deposits"`
	Channels    int64 `json:"channels"`
}
