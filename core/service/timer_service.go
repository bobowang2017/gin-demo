package service

import (
	"encoding/json"
	"gin-demo/core/dto"
	m "gin-demo/core/model"
	"gin-demo/infra/common"
	"gin-demo/infra/utils/log"
	"github.com/kirinlabs/HttpRequest"
	"github.com/robfig/cron/v3"
	"net/http"
	"strings"
	"time"
)

type TimerTask struct {
	Task *m.Task
}

func (t *TimerTask) BeforeDo() bool {
	var (
		svc  = NewTimerService()
		task *m.Task
		err  error
	)
	if task, err = svc.GetById(t.Task.ID); err != nil {
		log.Logger.Error(err)
		return false
	}
	stopAt := task.StopAt
	if !stopAt.IsZero() && stopAt.Before(time.Now()) {
		if t.Task.TaskId != 0 {
			common.TimerCron.Remove(cron.EntryID(task.TaskId))
			_ = svc.DeleteById(task.ID)
			return false
		}
	}
	return true
}

func (t *TimerTask) Do() {
	var (
		req       *HttpRequest.Request
		resp      *HttpRequest.Response
		paramsDto dto.TimerTaskDto
		err       error
	)
	if !t.BeforeDo() {
		return
	}
	log.Logger.Infof("定时任务%s开始执行", t.Task.Name)
	params := t.Task.Params
	if err = json.Unmarshal([]byte(params), &paramsDto); err != nil {
		log.Logger.Error("反序列化异常")
		return
	}
	req = HttpRequest.NewRequest().SetTimeout(3)
	if paramsDto.Header != nil {
		req.SetHeaders(paramsDto.Header)
	}
	switch strings.ToUpper(paramsDto.Method) {
	case "GET":
		resp, err = req.JSON().Get(paramsDto.Url)
	case "POST":
		resp, err = req.JSON().Post(paramsDto.Url, paramsDto.Args)
	default:
		resp, err = req.JSON().Get(paramsDto.Url, paramsDto.Args)
	}
	if err != nil {
		log.Logger.Error(err)
		return
	}
	if resp != nil {
		defer resp.Close()
	}
	if resp.StatusCode() != http.StatusOK {
		log.Logger.Error(err)
		return
	}
	log.Logger.Infof("定时任务%s执行成功", t.Task.Name)
}

type TimerTaskOpt func(*TimerTask)

func WithTask(task *m.Task) TimerTaskOpt {
	return func(t *TimerTask) {
		t.Task = task
	}
}

func NewTimerTask(opts ...TimerTaskOpt) *TimerTask {
	client := &TimerTask{}
	for _, opt := range opts {
		opt(client)
	}
	return client
}
