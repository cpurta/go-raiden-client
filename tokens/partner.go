package tokens

import "github.com/ethereum/go-ethereum/common"

type Partner struct {
	Address    common.Address `json:"partner_address"`
	ChannelURI string         `json:"channel"`
}
