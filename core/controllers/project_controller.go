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
	router.POST("", projectController.CreateProject)
}

func (p *ProjectController) List(c *gin.Context) {
	var (
		queryDto = dto.ProjectListQuery{}
		result   []*service.ProjectTreeNode
		total    int64
		err      error
	)

	c.Copy()

	if err = c.BindQuery(&queryDto); err != nil {
		log.Logger.Error(err)
		common.RespInputErrorJSON(c, err.Error())
	}
	if queryDto.Page <= 0 {
		queryDto.Page = 1
	}
	if queryDto.Size <= 0 {
		queryDto.Size = 10
	}
	if result, total, err = p.projectService.List(&queryDto); err != nil {
		log.Logger.Error(err)
		common.RespInputErrorJSON(c, err.Error())
	}
	common.RespSuccessPageJSON(c, result, total)
}

func (p *ProjectController) CreateProject(c *gin.Context) {
	var (
		createDto = dto.ProjectCreateDto{}
		err       error
	)
	if err = c.BindJSON(&createDto); err != nil {
		common.RespInputErrorJSON(c, err.Error())
		return
	}
	common.RespSuccessJSON(c, nil)
}
