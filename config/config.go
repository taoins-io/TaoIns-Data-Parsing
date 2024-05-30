package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
)

type TomlConfig struct {
	DB     DbConfig `toml:"db"`
	Log    Log      `toml:"log"`
	Chain  Chain    `toml:"chain"`
	Redis  Redis    `toml:"redis"`
	Task   TASK     `toml:"task"`
	Server Server   `toml:"server"`
}

type Log struct {
	Level         string `toml:"level"`
	LogFileDir    string `toml:"path"`
	AppName       string `toml:"name"`
	ErrorFileName string `toml:"error"`
	WarnFileName  string `toml:"warn"`
	InfoFileName  string `toml:"info"`
	DebugFileName string `toml:"debug"`
	MaxSize       int    `toml:"max_size"`
	MaxBackups    int    `toml:"max_backups"`
	MaxAge        int    `toml:"max_age"`
	Config        zap.Config
}

type DbConfig struct {
	UserName string `toml:"user_name"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Database string `toml:"database"`
}

type Chain struct {
	BeginBlock       int64    `toml:"begin_block"`
	SecondEveryBlock int64    `toml:"second_every_block"`
	HttpAddr         []string `toml:"http_addr"`
}

type Redis struct {
	Address string `toml:"address"`
	DB      int    `toml:"db"`
}

type TASK struct {
	ChainEventDataCron string `toml:"chain_event_data_cron"`
	ChainBlockCron     string `toml:"chain_block_cron"`
	ChainClientCron    string `toml:"chain_client_cron"`
	MonitorBlockCron   string `toml:"monitor_block_cron"`
}

type Server struct {
	Port string `toml:"port"`
}

var Config TomlConfig

func InitConfig(filePath string) (err error) {
	if _, err = toml.DecodeFile(filePath, &Config); err != nil {
		fmt.Printf("failed to init file: %v", err)
	}
	return
}
