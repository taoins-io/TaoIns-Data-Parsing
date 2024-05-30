package cron

import (
	"github.com/robfig/cron/v3"
	"strconv"
	"tao/config"
	"tao/consts"
	"tao/database/gredis"
	"tao/logger"
	"tao/util"
	"time"
)

type BlockJob struct {
}

func (j *BlockJob) Run() {
	lastBlockNumber()
}

func BlockCron() {
	if !gredis.Lock(consts.ChainBlockLock, time.Minute) {
		return
	}
	defer gredis.UnLock(consts.ChainBlockLock)
	c := cron.New()
	cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger))
	spec := config.Config.Task.ChainBlockCron
	c.AddJob(spec, cron.NewChain(util.SkipIfStillRunningOrLock("BlockCron"), cron.Recover(cron.VerbosePrintfLogger(logger.GetLogger()))).Then(&BlockJob{}))
	c.Start()
}

func lastBlockNumber() {
	taoBlock, err := util.GetLastBlock()
	if taoBlock.BlockNumber < 1 {
		logger.GetLogger().Errorf("lastBlockNumber error: %v", err)
		return
	}
	gredis.SetStringExpiration(consts.ChainBlock, strconv.FormatInt(taoBlock.BlockNumber, 10), 0)
}
