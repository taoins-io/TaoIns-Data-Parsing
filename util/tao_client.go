package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"strconv"
	"tao/config"
	"tao/consts"
	"tao/database/gredis"
	"tao/logger"
	vo2 "tao/vo"
	"time"
)

func GetLastBlock() (taoBlock vo2.TaoBlock, err error) {
	url := config.Config.Chain.HttpAddr[0] + "/latest_block"
	res, err := GetMethod(url)
	if err != nil || res.StatusCode != 200 {
		logger.GetLogger().Errorf("GetLastBlocks err: %v, res: %v", err, string(res.Content))
		return vo2.TaoBlock{}, err
	}
	taoData, err := getTaoResultData(res)
	if err != nil {
		logger.GetLogger().Errorf("GetLastBlocks error: %s, url: %s", err, url)
		return
	}
	if err = json.Unmarshal(taoData.Data, &taoBlock); err != nil {
		logger.GetLogger().Errorf("GetLastBlocks error: %s, data: %s", err, string(taoData.Data))
		return vo2.TaoBlock{}, err
	}
	datumTime := vo2.DatumTime{
		BlockNumber: taoBlock.BlockNumber,
		Timestamp:   time.UnixMilli(taoBlock.Time),
	}
	datumTimeByte, _ := json.Marshal(datumTime)
	gredis.SetStringExpiration(consts.BlockDatumTime, string(datumTimeByte), time.Minute*10)
	return
}

func GetBlockTransferByNumber(number int64) (taoTransfers []vo2.TaoTransfer, err error) {
	url := config.Config.Chain.HttpAddr[0] + "/transfers/" + strconv.FormatInt(number, 10)
	res, err := GetMethod(url)
	if err != nil || res.StatusCode != 200 {
		logger.GetLogger().Errorf("GetBlockEventByNumber err: %v, res: %v", err, string(res.Content))
		return
	}
	taoData, err := getTaoResultData(res)
	if err != nil || taoData.Code != 0 {
		logger.GetLogger().Errorf("GetBlockEventByNumber error: %s, url: %s", err, url)
		return taoTransfers, errors.New(fmt.Sprintf("Service error code:%v, message: %v", taoData.Code, string(taoData.Data)))
	}
	if err = json.Unmarshal(taoData.Data, &taoTransfers); err != nil {
		logger.GetLogger().Errorf("GetBlockEventByNumber error: %s, data: %s", err, string(taoData.Data))
		return
	}
	return
}

func getTaoResultData(res Res) (taoData vo2.TaoData, err error) {
	err = json.Unmarshal(res.Content, &taoData)
	if err != nil {
		logger.GetLogger().Errorf("GetLastBlocks UnmarshalJson err %v, data: %s", err, string(res.Content))
		return
	}
	return
}

func GetNewBlock() int64 {
	var blockNumber int64
	newBlockStr, err := gredis.GetValue(consts.ChainBlock)
	if err != nil {
		logger.GetLogger().Errorf("GetNewBlock GetValue err %v", zap.Error(err))
		return 0
	}
	blockNumber, _ = strconv.ParseInt(newBlockStr, 10, 64)
	if blockNumber < 1 {
		taoBlock, err := GetLastBlock()
		if err != nil {
			logger.GetLogger().Errorf("GetNewBlock error")
		} else {
			gredis.SetStringExpiration(consts.ChainBlock, strconv.FormatUint(uint64(taoBlock.BlockNumber), 10), 0)
		}
	}
	return blockNumber
}

//Get the corresponding block time based on the block

func TimeByHeight(number int64) time.Time {
	datumTimeStr, _ := gredis.GetValue(consts.BlockDatumTime)
	if len(datumTimeStr) < 1 {
		taoBlock, err := GetLastBlock()
		if err != nil {
			logger.GetLogger().Error("TimeByHeight: %v", err)
			return time.Now()
		}
		datumTime := vo2.DatumTime{
			BlockNumber: taoBlock.BlockNumber,
			Timestamp:   time.UnixMilli(taoBlock.Time),
		}
		datumTimeByte, _ := json.Marshal(datumTime)
		gredis.SetStringExpiration(consts.BlockDatumTime, string(datumTimeByte), time.Minute*10)
		return datumTime.Timestamp
	}
	var datumTime vo2.DatumTime
	if err := json.Unmarshal(bytes.NewBufferString(datumTimeStr).Bytes(), &datumTime); err != nil {
		taoBlock, err := GetLastBlock()
		if err != nil {
			logger.GetLogger().Error("TimeByHeight: %v", err)
			return time.Now()
		}
		datumTime.BlockNumber = taoBlock.BlockNumber
		datumTime.Timestamp = time.UnixMilli(taoBlock.Time)
		datumTimeByte, _ := json.Marshal(datumTime)
		gredis.SetStringExpiration(consts.BlockDatumTime, string(datumTimeByte), time.Minute*10)
		return datumTime.Timestamp
	}
	secondEveryBlock := decimal.NewFromInt(config.Config.Chain.SecondEveryBlock)
	secondSub := number - datumTime.BlockNumber
	blockTime := datumTime.Timestamp.Add(time.Millisecond * time.Duration(secondSub*secondEveryBlock.IntPart()))
	return blockTime
}
