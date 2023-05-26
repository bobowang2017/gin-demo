package controllers

import (
	"gin-demo/infra/common"
	"github.com/gin-gonic/gin"
	"time"
)

type SocketController struct {
}

func SocketControllerRegister(router *gin.RouterGroup) {
	socketController := SocketController{}
	router.GET("", socketController.Ping)
}

func (s *SocketController) Ping(c *gin.Context) {
	ws, err := common.SocketUpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.RespInternalErrorJSON(c, "WebSocket连接失败")
		return
	}
	defer func() {
		_ = ws.Close()
	}()
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-ticker.C:
			//err = ws.WriteMessage()
			if err != nil {
				break
			}
		}

		// 读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		// 写入ws数据
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}
