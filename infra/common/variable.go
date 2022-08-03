package common

import (
	"github.com/robfig/cron/v3"
)

//定义全局定时器
var TimerCron = cron.New()
