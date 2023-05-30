package controllers

import (
	"gin-demo/infra/common"
	infraSvc "gin-demo/infra/service"
	"github.com/gin-gonic/gin"
)

type SocketController struct {
}

func SocketControllerRegister(router *gin.RouterGroup) {
	socketController := SocketController{}
	router.GET("ws", socketController.WebSocket)
	router.GET("broadcast", socketController.Broadcast)
}

func (s *SocketController) WebSocket(c *gin.Context) {
	conn, err := common.WsUpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.RespInternalErrorJSON(c, "WebSocket连接失败")
		return
	}
	defer func() {
		_ = conn.Close()
	}()
	client := infraSvc.GetWsClient(infraSvc.WithConnOpt(conn))
	client.Manager.Register <- client
	go client.Write()
	client.Read()
}

func (s *SocketController) Broadcast(c *gin.Context) {
	infraSvc.GetWsClientManager().Broadcast <- []byte("broadcast")
	common.RespSuccessJSON(c, nil)
}
