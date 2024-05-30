package vo

import (
	"time"
)

type DatumTime struct {
	BlockNumber int64     `json:"blockNumber"`
	Timestamp   time.Time `json:"timestamp"`
}
