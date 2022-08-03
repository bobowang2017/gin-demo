package dao

import (
	"gin-demo/infra/model"
)

type SystemConfigDao struct {
}

func NewSystemConfigDao() *SystemConfigDao {
	return &SystemConfigDao{}
}

func (s *SystemConfigDao) GetUsingCfg() (*model.SystemConfig, error) {
	cfg := &model.SystemConfig{}
	if err := model.DB.Where(map[string]interface{}{"is_using": 1}).Find(cfg).Error; err != nil {
		return nil, err
	}
	return cfg, nil
}
