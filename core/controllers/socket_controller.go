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
	client := &infraSvc.WsClient{
		Manager:    infraSvc.GetWsClientManager(),
		SocketConn: conn,
		Send:       make(chan []byte, 256),
	}
	client.Manager.Register <- client
	go client.Read()
	go client.Write()
}
