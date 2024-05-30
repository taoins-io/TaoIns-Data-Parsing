package gdb

import (
	"tao/database/gdb/table"
	"tao/logger"
)

func (object *ChainDB) SaveTaoTransfer(transfer table.TaoTransfer) {
	err := object.db.Model(&table.TaoTransfer{}).Where("block_num = ? and event_index = ?", transfer.Block, transfer.EventIndex).FirstOrCreate(&transfer).Error
	if err != nil {
		logger.GetLogger().Errorf("SaveTaoTransfer err %v", err)
	}
}
