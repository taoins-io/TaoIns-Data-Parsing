package cmd

import "tao/cron"

func (p *Program) chainRun() {
	cron.BlockCron()
	cron.PullChainDataCron()
}
