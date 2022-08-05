package settings

import (
	"fmt"
	"gin-demo/infra/common"
	"gin-demo/infra/utils/config"
	"time"
)

type SettingConfig struct {
	System   config.System   `yaml:"System"`
	Server   config.Server   `yaml:"Server"`
	Database config.Database `yaml:"Database"`
	Redis    config.Redis    `yaml:"Redis"`
}

var Config = &SettingConfig{}

// Setup 正常程序启动加载配置
func Setup() {
	config.LoadConfig(fmt.Sprintf("%s-", common.ModuleCore), Config, false)
	Config.Server.ReadTimeout = Config.Server.ReadTimeout * time.Second
	Config.Server.WriteTimeout = Config.Server.WriteTimeout * time.Second
}

// SetupTest 测试类启动加载配置
func SetupTest() {
	config.LoadConfig(fmt.Sprintf("%s-", common.ModuleCore), Config, true)
	Config.Server.ReadTimeout = Config.Server.ReadTimeout * time.Second
	Config.Server.WriteTimeout = Config.Server.WriteTimeout * time.Second
}