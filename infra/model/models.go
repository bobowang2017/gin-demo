package model

import (
	"fmt"
	"gin-demo/core/settings"
	"gin-demo/infra/common"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
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
	var ormLogger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			// Ignore ErrRecordNotFound error for logger
			IgnoreRecordNotFoundError: true,
		})
	if gin.Mode() == "debug" {
		ormLogger.LogMode(logger.Info)
		//ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger.LogMode(logger.Warn)
		//ormLogger = logger.Default
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		settings.Config.Database.User,
		settings.Config.Database.Password,
		settings.Config.Database.Host,
		settings.Config.Database.Name)
	DB, err = gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			Logger:         ormLogger,
			NamingStrategy: schema.NamingStrategy{SingularTable: true, TablePrefix: settings.Config.Database.TablePrefix},
		},
	)

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	sqlDB.SetMaxIdleConns(settings.Config.Database.MaxIdleCon)
	sqlDB.SetMaxOpenConns(settings.Config.Database.MaxOpenCon)
	// 数据库连接最大生存时间(默认是8小时,也就是说如果8小时内没有任何数据库操作的话,数据库就会关闭连接,当前线上设置的是1800秒)
	sqlDB.SetConnMaxIdleTime(time.Duration(settings.Config.Database.MaxIdleTime) * time.Second)
}
