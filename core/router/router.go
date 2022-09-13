package router

import (
	c "gin-demo/core/controllers"
	"gin-demo/infra/middlewares"
	"gin-demo/infra/validators"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitRouter() *gin.Engine {

	// 自定义验证器的注册
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("not_allow_blank", validators.NotAllowBlank); err != nil {
			panic(err)
		}
	}

	router := gin.New()
	// 要在路由组之前全局使用「跨域中间件」, 否则OPTIONS会返回404
	router.Use(middlewares.Cors(), middlewares.Auth(), middlewares.Recover)
	timerGroup := router.Group("/api/v1/timers")
	{
		c.TimerRouterRegister(timerGroup)
	}

	projectGroup := router.Group("/api/v1/projects")
	{
		c.ProjectControllerRegister(projectGroup)
	}

	router.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "200",
			"msg":    "success",
			"data":   nil,
		})
	})
	return router
}
