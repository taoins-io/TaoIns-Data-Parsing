package event_factory

import (
	"strconv"
	"tao/database/gdb"
	"tao/database/gdb/table"
	"tao/util"
	"tao/vo"
)

var contractTemplate = make(map[string]bool)

const (
	EventType = "ffffffff"
)

func EventAllFactory(eventNode vo.EventNode) {
	contractOnce.Do(func() {
		contractTemplate[EventType] = true
	})
	// check inscription
	//ffffffff13000000000000692e74616f7562692e636f6d2f616972642e706e67
	assetType := eventNode.ToHex[8:9]
	contentType := eventNode.ToHex[9:10]
	taoTransfer := table.TaoTransfer{
		Sender:      eventNode.From,
		Receiver:    eventNode.To,
		ReceiverHex: eventNode.ToHex,
		Amount:      eventNode.Amount,
		Block:       eventNode.BlockNumber,
		BlockTime:   util.TimeByHeight(eventNode.BlockNumber).UnixMilli(),
		EventIndex:  strconv.FormatInt(eventNode.Id, 10),
		AssetType:   assetType,
		ContentType: contentType,
	}
	gdb.Inst().SaveTaoTransfer(taoTransfer)
}
