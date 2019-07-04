package pendingtransfers

import "github.com/ethereum/go-ethereum/common"

type Transfer struct {
	ChannelIdentifier      int64          `json:"channel_identifier"`
	Initiator              common.Address `json:"initiator"`
	LockedAmount           int64          `json:"locked_amount"`
	PaymentIdentifier      int64          `json:"payment_identifier"`
	Role                   string         `json:"role"`
	Target                 common.Address `json:"target"`
	TokenAddress           common.Address `json:"token_address"`
	TokenNetworkIdentifier common.Address `json:"token_network_identifier"`
	TransferredAmount      int64          `json:"transferred_amount"`
}
