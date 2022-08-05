package model

import (
	"gin-demo/infra/common"
	m "gin-demo/infra/model"
)

type Task struct {
	m.BaseModel
	Name        string          `gorm:"not null;" json:"name"`
	Description string          `gorm:"default:null;" json:"description"`
	Cron        string          `gorm:"not null;" json:"cron"`
	Params      string          `gorm:"not null;" json:"params"`
	TaskId      int             `gorm:"default:null;" json:"taskId"`
	StopAt      common.JSONTime `gorm:"default:null;" json:"stopAt"`
}
