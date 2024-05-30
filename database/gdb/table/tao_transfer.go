package table

type TaoTransfer struct {
	ID          int64  `gorm:"primarykey;column:id"`
	Updated     int64  `gorm:"column:updated;autoUpdateTime:milli"`
	Created     int64  `gorm:"column:created;autoCreateTime:milli"`
	AssetType   string `gorm:"column:asset_type;"`
	ContentType string `gorm:"column:content_type;"`
	Sender      string `gorm:"column:sender;"`
	Receiver    string `gorm:"column:receiver;"`
	ReceiverHex string `gorm:"column:receiver_hex;"`
	EventIndex  string `gorm:"column:event_index;"`
	Amount      int64  `gorm:"column:amount;"`
	Block       int64  `gorm:"column:block_num;"`
	BlockTime   int64  `gorm:"column:block_time;"`
}

func (TaoTransfer) TableName() string {
	return "tao_transfer"
}
