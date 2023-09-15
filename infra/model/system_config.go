package model

// SystemConfig   系统配置对象Model信息
type SystemConfig struct {
	BaseModel
	Name    string `gorm:"size:45;not null" json:"name"`
	Content string `gorm:"type:text;" json:"content"`
	IsUsing *int   `gorm:"not null;" json:"isUsing"`
}

// SysCfg 定义系统配置Content结构对象
type SysCfg struct {
	ValidCodeExpire int
	UserExpire      int
	NoAuthUrl       map[string]string
	TokenExpire     int
	UserLabelExpire int
}
