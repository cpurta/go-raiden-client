package address

import "github.com/ethereum/go-ethereum/common"

// Address is a wrapper arroung the ethereum 20-byte address.
type Address struct {
	Address common.Address `json:"address"`
}
