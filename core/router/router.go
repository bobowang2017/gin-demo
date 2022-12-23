package router

import (
	c "gin-demo/core/controllers"
	"gin-demo/infra/middlewares"
	"github.com/arl/statsviz"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	// 集成statsviz监控
	router.GET("/debug/statsviz/*filepath", func(context *gin.Context) {
		if context.Param("filepath") == "/ws" {
			statsviz.Ws(context.Writer, context.Request)
			return
		}
		statsviz.IndexAtRoot("/debug/statsviz").ServeHTTP(context.Writer, context.Request)
	})

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
