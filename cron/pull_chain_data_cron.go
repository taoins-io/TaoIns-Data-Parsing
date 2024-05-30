package cron

import (
	"github.com/robfig/cron/v3"
	"strconv"
	"tao/config"
	"tao/consts"
	"tao/database/gredis"
	"tao/event_factory"
	"tao/logger"
	"tao/util"
)

var eventFactory *event_factory.EventFactory

type PullChainDataJob struct {
}

func (j *PullChainDataJob) Run() {
	// Synchronize data
	eventFactory.FnPublicFuncStartHistoryAll()
}

func PullChainDataCron() *cron.Cron {
	pullBlock, _ := gredis.GetValue(consts.ChainEventBlock)
	if len(pullBlock) < 1 {
		gredis.SetValueExpiration(consts.ChainEventBlock, uint64(config.Config.Chain.BeginBlock), 0)
	} else {
		block, err := strconv.ParseInt(pullBlock, 10, 64)
		if err != nil {
			logger.GetLogger().Error(err)
		} else if block < config.Config.Chain.BeginBlock {
			gredis.SetValueExpiration(consts.ChainEventBlock, uint64(config.Config.Chain.BeginBlock), 0)
		}
	}
	//Pull data regularly
	c := cron.New()
	spec := config.Config.Task.ChainEventDataCron
	c.AddJob(spec, cron.NewChain(util.SkipIfStillRunningOrLock("PullChainDataCron"), cron.Recover(cron.VerbosePrintfLogger(logger.GetLogger()))).
		Then(&PullChainDataJob{}))
	c.Start()
	return c
}
