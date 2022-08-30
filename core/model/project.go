package model

import (
	m "gin-demo/infra/model"
)

type Project struct {
	m.BaseModel
	Name         string `gorm:"not null;" json:"name"`
	Code         string `gorm:"not null;" json:"code"`
	Age          int    `gorm:"not null;" json:"age"`
	ParentId     int    `gorm:"not null;" json:"parentId"`
	Level        int    `gorm:"not null;" json:"level"`
	RootParentId int    `gorm:"not null;" json:"rootParentId"`
	PathIds      string `gorm:"not null;" json:"pathIds"`
}
