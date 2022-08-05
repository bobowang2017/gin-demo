package controllers

import (
	"fmt"
	"gin-demo/core/dto"
	"gin-demo/core/service"
	"gin-demo/infra/common"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type TimerController struct {
	timerService service.ITimerService
}

func TimerRouterRegister(router *gin.RouterGroup) {
	timerController := TimerController{
		service.NewTimerService(),
	}
	router.GET("", timerController.List)
	router.POST("", timerController.Register)
}

func (t *TimerController) Register(c *gin.Context) {
	var (
		timerAddDto = &dto.TimerTaskAddDto{}
		err         error
	)
	if err = c.BindJSON(timerAddDto); err != nil {
		common.RespInputErrorJSON(c, err.Error())
		return
	}
	if err = t.timerService.RegisterTimer(timerAddDto); err != nil {
		common.RespInternalErrorJSON(c, err.Error())
		return
	}
	common.RespSuccessJSON(c, nil)
}

func (t *TimerController) List(c *gin.Context) {
	var (
		entryIds []interface{}
		result   []cron.Entry
	)
	result = common.TimerCron.Entries()
	for _, res := range result {
		fmt.Println(res.ID)
		entryIds = append(entryIds, res.ID)
	}
	common.RespSuccessJSON(c, entryIds)
}
