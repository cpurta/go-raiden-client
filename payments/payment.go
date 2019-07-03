package payments

import "github.com/ethereum/go-ethereum/common"

type Payment struct {
	InitiatorAddress common.Address `json:"initiator_address"`
	TargetAddress    common.Address `json:"target_address"`
	TokenAddress     common.Address `json:"token_address"`
	Amount           int64          `json:"amount"`
	Identifier       int64          `json:"identifier"`
}
