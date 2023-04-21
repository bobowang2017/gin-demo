package dao

import (
	m "gin-demo/core/model"
	"gin-demo/infra/dao"
)

type TimerDao struct {
	dao.BaseDao
}

func NewTimerDao() *TimerDao {
	return &TimerDao{}
}

func (t *TimerDao) GetAll() ([]*m.Task, error) {
	var (
		tasks []*m.Task
		err   error
	)
	err = t.BaseDao.GetObjByCondition(&m.Task{}, map[string]interface{}{}, &tasks)
	return tasks, err
}
