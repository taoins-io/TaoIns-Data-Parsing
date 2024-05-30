package cmd

import (
	"github.com/kardianos/service"
	"tao/database/gdb"
	"tao/database/gredis"
	"tao/logger"
)

type Program struct {
}

func (p *Program) Start(s service.Service) error {
	p.run()
	logger.GetLogger().Info("service start")
	return nil
}
func (p *Program) Stop(s service.Service) error {
	logger.GetLogger().Info("service stop")
	return nil
}
func (p *Program) run() {
	gdb.Inst()
	err := gredis.InitClient()
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	p.chainRun()
}
