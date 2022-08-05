package model

import (
	"fmt"
	s "gin-demo/core/settings"
	"gin-demo/infra/common"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
)

var DB *gorm.DB

type BaseModel struct {
	ID          int             `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedTime common.JSONTime `gorm:"autoCreateTime" json:"createdTime"`
	UpdatedTime common.JSONTime `gorm:"autoUpdateTime:milli" json:"updatedTime"`
}

type TableNameAble interface {
	TableName() string
}

// Setup initializes the database instance
func Setup() {
	var err error
	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local",
		s.Config.Database.User,
		s.Config.Database.Password,
		s.Config.Database.Host,
		s.Config.Database.Name)
	DB, err = gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			Logger:         ormLogger,
			NamingStrategy: schema.NamingStrategy{SingularTable: true, TablePrefix: s.Config.Database.TablePrefix},
		},
	)

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
}
