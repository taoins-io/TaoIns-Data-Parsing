package gdb

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logDb "gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"tao/config"
	"tao/logger"
	"time"
)

var (
	dbInst *ChainDB
	dbOnce sync.Once
)

type ChainDB struct {
	db *gorm.DB
	sync.Mutex
}

func newDB() *ChainDB {
	object := &ChainDB{}
	args := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		config.Config.DB.UserName,
		config.Config.DB.Password,
		config.Config.DB.Host,
		config.Config.DB.Port,
		config.Config.DB.Database,
	)
	logLevel := logDb.Error
	if config.Config.Log.Level == "debug" {
		logLevel = logDb.Info
	}
	blogdown := logDb.New(logger.NewWriter(log.New(os.Stdout, "", log.LstdFlags)), logDb.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logLevel,
		Colorful:      false,
	})
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   blogdown,
	})
	db.Set("logger", logger.GetLogger())
	if err != nil {
		logger.GetLogger().Errorf("cannot connect to pg %v", err)
	}
	dbConn, err := db.DB()
	if err != nil {
		logger.GetLogger().Errorf("cannot get pg %v", err)
	}
	dbConn.SetMaxIdleConns(100)
	dbConn.SetMaxOpenConns(100)
	dbConn.SetConnMaxLifetime(time.Second * 200)
	object.db = db
	return object
}

func Inst() *ChainDB {
	dbOnce.Do(func() {
		dbInst = newDB()
	})
	return dbInst
}
