package cache

import (
	"github.com/dy-gopkg/kit/dao/redis"
	"github.com/dy-platform/user-srv-session/util"
)

func CacheInit() {
	redis.Init(util.DefaultRedisConf.Addr, util.DefaultRedisConf.Password,
		util.DefaultRedisConf.MaxIdle, util.DefaultRedisConf.MaxActive)
}
