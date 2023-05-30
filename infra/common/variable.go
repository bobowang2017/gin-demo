package common

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/robfig/cron/v3"
)

var (
	// TimerCron 定义全局定时器
	TimerCron = cron.New()
	// ValidTrans 定义错误翻译对象
	ValidTrans ut.Translator
	ValidObj   *validator.Validate
)

var WsUpGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
