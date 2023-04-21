package service

import (
	"encoding/json"
	"gin-demo/core/dao"
	"gin-demo/core/dto"
	m "gin-demo/core/model"
	"gin-demo/infra/common"
	"gin-demo/infra/utils"
	"gin-demo/infra/utils/log"
)

type ITimerService interface {
	GetAll() ([]*m.Task, error)
	UpdateById(timerId int, params map[string]interface{}) error
	DeleteById(timerId int) error
	GetById(id int) (*m.Task, error)
	RegisterTimer(taskAddDto *dto.TimerTaskAddDto) error
}

// 关联模块接口实现
type timerService struct {
	timerDao *dao.TimerDao
}

func (t *timerService) GetById(id int) (*m.Task, error) {
	var (
		result = &m.Task{}
		err    error
	)
	err = t.timerDao.GetObjById(result, id)
	return result, err
}

func (t *timerService) GetAll() ([]*m.Task, error) {
	return t.timerDao.GetAll()
}

func (t *timerService) UpdateById(timerId int, params map[string]interface{}) error {
	return t.timerDao.UpdateObjById(&m.Task{}, timerId, params)
}

func (t *timerService) DeleteById(timerId int) error {
	return t.timerDao.DeleteObjById(&m.Task{}, timerId)
}

func (t *timerService) RegisterTimer(taskAddDto *dto.TimerTaskAddDto) error {
	var (
		taskId    interface{}
		taskModel *m.Task
		err       error
	)
	paramStr, _ := json.Marshal(taskAddDto.Params)
	taskModel = &m.Task{
		Name:        taskAddDto.Name,
		Description: taskAddDto.Description,
		Cron:        taskAddDto.Cron,
		Params:      string(paramStr),
		StopAt:      common.JSONTime{Time: utils.StringToTime(taskAddDto.StopAt)},
	}
	if err = t.timerDao.Create(taskModel); err != nil {
		log.Logger.Error(err)
		return err
	}

	timerTask := NewTimerTask(WithTask(taskModel))
	if taskId, err = common.TimerCron.AddFunc(taskModel.Cron, timerTask.Do); err != nil {
		log.Logger.Errorf("任务注册失败|name=%s|%s", taskModel.Name, err.Error())
		return err
	}
	if err = t.UpdateById(taskModel.ID, map[string]interface{}{"task_id": taskId}); err != nil {
		log.Logger.Errorf("更新任务TaskId失败|name=%s", taskModel.Name)
		_ = t.DeleteById(taskModel.ID)
		return err
	}
	return nil
}

func NewTimerService() ITimerService {
	return &timerService{dao.NewTimerDao()}
}
