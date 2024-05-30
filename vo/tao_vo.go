package vo

import (
	"encoding/json"
)

type TaoData struct {
	Code int16           `json:"code"`
	Data json.RawMessage `json:"data"`
}

type TaoBlock struct {
	BlockNumber int64 `json:"blockNumber"`
	Time        int64 `json:"timestamp"`
}

type TaoTransfer struct {
	EventIndex  int64           `json:"id"`
	BlockNumber string          `json:"blockHeight"`
	Data        json.RawMessage `json:"data"`
}

type EventNode struct {
	BlockNumber int64  `json:"blockNumber"`
	Id          int64  `json:"id"`
	From        string `json:"from"`
	To          string `json:"to"`
	ToHex       string `json:"to_hex"`
	Amount      int64  `json:"amount"`
}
