package dao

import (
	m "gin-demo/core/model"
	"gin-demo/infra/dao"
	"gin-demo/infra/model"
)

type ProjectDao struct {
	dao.BaseDao
}

func NewProjectDao() *ProjectDao {
	return &ProjectDao{}
}

func (p *ProjectDao) GetRootParentIdByCondition(page, size int) ([]int, error) {
	var (
		ids = make([]int, 0)
		err error
	)
	err = model.DB.Model(&m.Project{}).Select("distinct(root_parent_id)").Where(
		"parent_id = ?", -1).Limit(size).Offset((page - 1) * size).Scan(&ids).Error
	return ids, err
}

func (p *ProjectDao) GetTotalDistinctRootParentId() (int64, error) {
	var (
		total int64
		err   error
	)
	err = model.DB.Model(&m.Project{}).Select("count(distinct(root_parent_id)) as total").Where(
		"parent_id = ?", -1).Scan(&total).Error
	return total, err
}

func (p *ProjectDao) GetByRootParentId(rootParentIds []int) ([]*m.Project, error) {
	var (
		rows = make([]*m.Project, 0)
		err  error
	)
	err = model.DB.Model(&m.Project{}).Where(
		"root_parent_id in ?", rootParentIds).Scan(&rows).Error
	return rows, err
}
