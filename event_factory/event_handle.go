package event_factory

import (
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"sort"
	"strconv"
	"sync"
	"tao/consts"
	"tao/database/gredis"
	"tao/logger"
	"tao/util"
	"tao/vo"
)

var contractOnce sync.Once

type EventFactory struct {
	AllContractTypeMap map[string]int8
}

var WsQuit = make(chan int)

func (e *EventFactory) FnPublicFuncStartHistoryAll() error {
	var err error
	pulledBlock, _ := gredis.GetValue(consts.ChainEventBlock)
	if len(pulledBlock) < 1 {
		logger.GetLogger().Errorf("pulledBlock Err:%+v", err)
		return err
	}
	block, _ := strconv.ParseInt(pulledBlock, 10, 64)
	block += 1
	newBlockNumber := util.GetNewBlock()
	if newBlockNumber < 1 {
		return nil
	}
	if block >= newBlockNumber {
		return nil
	}
	taoTransfers, err := util.GetBlockTransferByNumber(block)
	if err != nil {
		logger.GetLogger().Errorf("GetBlockEventByNumber[%v] Err:%v", block, err)
		return err
	}
	// sort by EventIndex
	sort.SliceStable(taoTransfers, func(i, j int) bool {
		return taoTransfers[i].EventIndex < taoTransfers[j].EventIndex
	})
	for _, taoTransfer := range taoTransfers {
		e.Process(taoTransfer)
	}
	gredis.SetValueExpiration(consts.ChainEventBlock, uint64(block), 0)
	logger.GetLogger().Infof("Pull [%v] End", block)
	return nil
}

func (e *EventFactory) Process(taoTransfer vo.TaoTransfer) {
	var detail []interface{}
	if err := json.Unmarshal(taoTransfer.Data, &detail); err != nil {
		logger.GetLogger().Errorf("json.Unmarshal Err:%+v", err)
		return
	}
	toByte := base58.Decode(detail[1].(string))
	toHex := common.Bytes2Hex(toByte[:])
	toHex = toHex[2 : len(toHex)-4]
	blockNumber, _ := strconv.ParseInt(taoTransfer.BlockNumber, 10, 64)
	amount := detail[2].(float64)
	EventAllFactory(vo.EventNode{
		BlockNumber: blockNumber,
		Id:          taoTransfer.EventIndex,
		ToHex:       toHex,
		From:        detail[0].(string),
		To:          detail[1].(string),
		Amount:      int64(amount),
	})
}
