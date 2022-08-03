package controllers

import (
	"gin-demo/infra/common"
	"github.com/gin-gonic/gin"
)

type DemoController struct {
}

func DemoControllerRouterRegister(router *gin.RouterGroup) {
	demoController := DemoController{

	}
	router.GET("", demoController.Register)
}

func (t *DemoController) Register(c *gin.Context) {

	common.RespSuccessJSON(c, nil)
}
