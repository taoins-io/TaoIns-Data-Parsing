package gredis

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v8" // 注意导入的是新版本
	"log"
	"sync"
	"tao/config"
	"tao/consts"
	"time"
)

var (
	Rdb   *redis.ClusterClient
	mutex sync.Mutex
)

// init connect

func InitClient() (err error) {
	Rdb = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{config.Config.Redis.Address},
		Password: "",
		PoolSize: 100,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err = Rdb.Ping(ctx).Result()
	return err
}

func SetStringExpiration(key string, val string, duration time.Duration) error {
	ctx := context.Background()
	return Rdb.Set(ctx, key, val, duration).Err()
}

func SetValueExpiration(key string, val uint64, duration time.Duration) error {
	ctx := context.Background()
	return Rdb.Set(ctx, key, val, duration).Err()
}

func GetValue(key string) (string, error) {
	ctx := context.Background()
	return Rdb.Get(ctx, key).Result()
}

func Delete(key string) error {
	ctx := context.Background()
	return Rdb.Del(ctx, key).Err()
}

func Lock(key string, expiration time.Duration) bool {
	mutex.Lock()
	defer mutex.Unlock()
	key = fmt.Sprintf(consts.RedisLock, key)
	result, err := Rdb.SetNX(context.Background(), key, 1, expiration).Result()
	if err != nil {
		log.Println(err.Error())
	}
	return result
}

func UnLock(key string) {
	key = fmt.Sprintf(consts.RedisLock, key)
	err := Delete(key)
	if err != nil {
		log.Println(err.Error())
	}
}
