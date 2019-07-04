package payments

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type event struct {
	EventName  string `json:"event"`
	Amount     int64  `json:"amount"`
	Initiator  string `json:"initiator"`
	Target     string `json:"target"`
	Identifier int64  `json:"identifier"`
	LogTime    string `json:"log_time"`
}

type Event struct {
	EventName  string
	Amount     int64
	Initiator  common.Address
	Target     common.Address
	Identifier int64
	LogTime    time.Time
}
