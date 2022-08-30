package controllers

import (
	"gin-demo/core/dto"
	"gin-demo/core/service"
	"gin-demo/infra/common"
	"gin-demo/infra/utils/log"
	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	projectService service.IProjectService
}

func ProjectControllerRegister(router *gin.RouterGroup) {
	projectController := ProjectController{
		service.NewProjectService(),
	}
	router.GET("", projectController.List)
}

func (p *ProjectController) List(c *gin.Context) {
	var (
		queryDto = dto.ProjectListQuery{}
		result   []*service.ProjectTreeNode
		err      error
	)

	if err = c.BindQuery(&queryDto); err != nil {
		log.Logger.Error(err)
		common.RespInputErrorJSON(c, err.Error())
	}
	if queryDto.Page == 0 {
		queryDto.Page = 1
	}
	if queryDto.Size == 0 {
		queryDto.Size = 10
	}
	if result, err = p.projectService.List(&queryDto); err != nil {
		log.Logger.Error(err)
		common.RespInputErrorJSON(c, err.Error())
	}
	common.RespSuccessJSON(c, result)
}