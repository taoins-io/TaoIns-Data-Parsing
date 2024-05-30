package util

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"tao/database/gredis"
	"tao/logger"
	"time"
)

// SkipIfStillRunningOrLock

func SkipIfStillRunningOrLock(jobFuncName string) cron.JobWrapper {
	var ch = make(chan struct{}, 1)
	ch <- struct{}{}
	return func(j cron.Job) cron.Job {
		return cron.FuncJob(func() {
			if !gredis.Lock(jobFuncName, 5*time.Minute) {
				return
			}
			defer func() {
				gredis.UnLock(jobFuncName)
			}()
			select {
			case v := <-ch:
				j.Run()
				ch <- v
			default:
				logger.GetLogger().Warnw("skip job", zap.Any("jobName", jobFuncName))
			}
		})
	}
}
