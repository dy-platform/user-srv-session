package main

import (
	"github.com/dy-gopkg/kit/micro"
	"github.com/dy-platform/user-srv-session/cache"
	"github.com/dy-platform/user-srv-session/handler"
	"github.com/dy-platform/user-srv-session/idl/platform/user/srv-session"
	"github.com/dy-platform/user-srv-session/util"
	"github.com/sirupsen/logrus"
)

func main() {
	micro.Init()

	// 初始化配置
	util.Init()

	//TODO 初始化缓存
	cache.CacheInit()

	err := platform_user_srv_session.RegisterSessionHandler(micro.Server(), &handler.Handler{})
	if err != nil {
		logrus.Fatalf("RegisterPassportHandler error:%v", err)
	}

	micro.Run()
}
