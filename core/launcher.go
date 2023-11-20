package core

import (
	"gin-demo/core/settings"
	"gin-demo/infra/model"
	infraSvc "gin-demo/infra/service"
	"gin-demo/infra/utils/log"
	"gin-demo/infra/utils/mq"
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
	// WebSocket管理服务启动
	go infraSvc.GetWsClientManager().Start()
	gin.SetMode(settings.Config.Server.RunMode)
	routers := InitRouter()
	if err := routers.Run("0.0.0.0:" + settings.Config.Server.HttpPort); err != nil {
		log.Logger.Fatalf("服务启动失败: %v", err)
	}
}

func StartConsumer() {
	setUp()
	mq.StartConsumer(nil)
}
