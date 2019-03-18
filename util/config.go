package util

import (
	"github.com/micro/go-config"
	"github.com/sirupsen/logrus"
)

type (
	RedisConfig struct {
		Addr      string `yaml:"addr"`
		Username  string `yaml:"username"`
		Password  string `yaml:"password"`
		MaxIdle   int    `yaml:"maxIdle"`
		MaxActive int    `yaml:"maxActive"`
	}
)

var (
	DefaultRedisConf RedisConfig
)

func Init() {
	//加载Redis 配置
	err := config.Get("Redis").Scan(&DefaultRedisConf)
	if err != nil {
		logrus.Fatalf("get Redis config error: %s", err)
	}

	if len(DefaultRedisConf.Addr) == 0 {
		logrus.Fatalf("invalid Redis addr")
	}

}
