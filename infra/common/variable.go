package common

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
)

// TimerCron 定义全局定时器
var TimerCron = cron.New()

var ValidTrans ut.Translator
var ValidObj *validator.Validate
