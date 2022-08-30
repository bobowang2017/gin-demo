package service

import (
	"fmt"
	"gin-demo/core/dao"
	"gin-demo/core/dto"
	m "gin-demo/core/model"
)

// ProjectTreeNode 定义Project树状结构体
type ProjectTreeNode struct {
	*m.Project
	Level    int                `json:"level"`
	Children []*ProjectTreeNode `json:"children,omitempty"`
}

type IProjectService interface {
	List(queryDto *dto.ProjectListQuery) ([]*ProjectTreeNode, error)
}

type projectService struct {
	projectDao *dao.ProjectDao
}

func (p *projectService) List(queryDto *dto.ProjectListQuery) ([]*ProjectTreeNode, error) {
	var (
		rows []*m.Project
		err  error
	)
	if queryDto.Age == 0 && queryDto.Name == "" && queryDto.Code == "" {
		rows, err = p.listWithCondition(queryDto.Page, queryDto.Size)
	} else {
		rows, err = p.listWithCondition(queryDto.Page, queryDto.Size)
	}
	fmt.Println(rows)
	return nil, err
}

func (p *projectService) listWithOutCondition(page, size int) ([]*m.Project, error) {
	var (
		result  []*m.Project
		rootIds []int
		err     error
	)
	if rootIds, err = p.projectDao.GetRootParentIdByCondition(page, size); err != nil {
		return result, err
	}
	if result, err = p.projectDao.GetByRootParentId(rootIds); err != nil {
		return result, err
	}
	return nil, nil
}

func (p *projectService) generateTree(rows []*m.Project) []*ProjectTreeNode {
	return nil
}

func (p *projectService) listWithCondition(page, size int) ([]*m.Project, error) {

	return nil, nil
}

func NewProjectService() IProjectService {
	return &projectService{dao.NewProjectDao()}
}
