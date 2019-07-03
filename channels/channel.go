package channels

import "github.com/ethereum/go-ethereum/common"

type channel struct {
	TokenNetworkIdentifier string `json:"token_network_identifier"`
	ChannelIdentifier      int64  `json:"channel_identifier"`
	PartnerAddress         string `json:"partner_address"`
	TokenAddress           string `json:"token_address"`
	Balance                int64  `json:"balance"`
	TotalDeposit           int64  `json:"total_deposit"`
	State                  string `json:"state"`
	SettleTimeout          int64  `json:"settle_timeout"`
	RevealTimeout          int64  `json:"reveal_timeout"`
}

type Channel struct {
	TokenNetworkIdentifier common.Address
	ChannelIdentifier      int64
	PartnerAddress         common.Address
	TokenAddress           common.Address
	Balance                int64
	TotalDeposit           int64
	State                  string
	SettleTimeout          int64
	RevealTimeout          int64
}
