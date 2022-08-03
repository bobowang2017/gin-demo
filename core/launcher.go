package core

import (
	"gin-demo/core/router"
	"gin-demo/core/settings"
	"gin-demo/infra/model"
	"gin-demo/infra/utils/log"
	"github.com/gin-gonic/gin"
)

func setUp() {
	settings.Setup()
	//redis.SetUp(settings.Config.Redis)
	log.Setup()
	model.Setup()
}

func Start() {
	setUp()
	gin.SetMode(settings.Config.Server.RunMode)
	routers := router.InitRouter()
	if err := routers.Run("0.0.0.0:" + settings.Config.Server.HttpPort); err != nil {
		log.Logger.Fatalf("服务启动失败: %v", err)
	}
}
