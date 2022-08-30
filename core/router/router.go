package router

import (
	c "gin-demo/core/controllers"
	"gin-demo/infra/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
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
