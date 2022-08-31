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

func (p *ProjectDao) GetRootParentIdByCondition(page, size int, params map[string]interface{}) ([]int, error) {
	var (
		ids = make([]int, 0)
		err error
	)
	err = model.DB.Model(&m.Project{}).Select("distinct(root_parent_id)").Where(
		"parent_id = ?", -1).Where(params).Limit(size).Offset((page - 1) * size).Scan(&ids).Error
	return ids, err
}

func (p *ProjectDao) GetDisRootParentIds(page, size int, params map[string]interface{}) ([]int, error) {
	var (
		ids = make([]int, 0)
		err error
	)

	db := model.DB.Model(&m.Project{}).Select("distinct(root_parent_id)")

	if v, ok := params["name"]; ok {
		db = db.Where("name like ?", "%"+v.(string)+"%")
	}
	if v, ok := params["age"]; ok {
		db = db.Where("age = ?", v)
	}
	if v, ok := params["code"]; ok {
		db = db.Where("code = ?", v)
	}

	err = db.Limit(size).Offset((page - 1) * size).Scan(&ids).Error
	return ids, err
}

func (p *ProjectDao) GetTotalDisRootParentId(params map[string]interface{}) (int64, error) {
	var (
		total int64
		err   error
	)

	db := model.DB.Model(&m.Project{}).Select("count(distinct(root_parent_id)) as total")

	if v, ok := params["name"]; ok {
		db = db.Where("name like ?", "%"+v.(string)+"%")
	}
	if v, ok := params["age"]; ok {
		db = db.Where("age = ?", v)
	}
	if v, ok := params["code"]; ok {
		db = db.Where("code = ?", v)
	}
	err = db.Scan(&total).Error
	return total, err
}

func (p *ProjectDao) GetPathIdsByRootParentIds(rootParentIds []int, params map[string]interface{}) ([]string, error) {
	var (
		pathIds = make([]string, 0)
		err     error
	)

	db := model.DB.Model(&m.Project{}).Select("path_ids")

	if v, ok := params["name"]; ok {
		db = db.Where("name like ?", "%"+v.(string)+"%")
	}
	if v, ok := params["age"]; ok {
		db = db.Where("age = ?", v)
	}
	if v, ok := params["code"]; ok {
		db = db.Where("code = ?", v)
	}

	err = db.Where("root_parent_id in ?", rootParentIds).Scan(&pathIds).Error
	return pathIds, err
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

func (p *ProjectDao) GetByRootParentIds(rootParentIds []int) ([]*m.Project, error) {
	var (
		rows = make([]*m.Project, 0)
		err  error
	)
	err = model.DB.Model(&m.Project{}).Where(
		"root_parent_id in ?", rootParentIds).Scan(&rows).Error
	return rows, err
}

func (p *ProjectDao) GetByIds(ids []int64) ([]*m.Project, error) {
	var (
		rows = make([]*m.Project, 0)
		err  error
	)
	err = model.DB.Model(&m.Project{}).Where("id in ?", ids).Scan(&rows).Error
	return rows, err
}
