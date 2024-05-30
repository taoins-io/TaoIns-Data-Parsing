package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tao/cmd"
	"tao/config"
	"tao/consts"
	logg "tao/logger"
	"time"
)

func main() {
	config.InitConfig(consts.ConfigPath)
	logg.InitLogger()
	go RunServer()
	cfg := &service.Config{
		Name:        "TaoInscription",
		DisplayName: "TaoInscription service",
		Description: "This is TaoInscription Go service.",
	}
	prg := &cmd.Program{}
	s, err := service.New(prg, cfg)
	if err != nil {
		logg.GetLogger().Errorf(err.Error())
	}
	logger, err := s.Logger(nil)
	if err != nil {
		logg.GetLogger().Errorf(err.Error())
	}
	if len(os.Args) == 2 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			logg.GetLogger().Errorf(err.Error())
		}
	} else {
		err = s.Run()
		if err != nil {
			logger.Error(err)
		}
	}
	if err != nil {
		logger.Error(err)
	}
}

func RunServer() {
	r := setupRouter()
	srv := &http.Server{
		Addr:    config.Config.Server.Port,
		Handler: r,
	}
	logg.GetLogger().Infof("server start port --> %v", srv.Addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.GetLogger().Errorf("listen err:%v", err)
		}
	}()
	// Wait for an interrupt signal to shut down the server gracefully (set a timeout of 5 seconds)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logg.GetLogger().Infof("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logg.GetLogger().Errorf("Shutdown err:%v", err)
	}
	logg.GetLogger().Infof("Server exiting")
}

func setupRouter() *gin.Engine {
	router := gin.New()
	return router
}
