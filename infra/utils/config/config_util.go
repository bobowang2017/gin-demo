package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
)

// 从制定文件名加载配置信息
// 默认文件路径"config/{filePrefix + filename}.yml"
func LoadConfig(filePrefix string, config interface{}, isTest bool) {
	var (
		currentPath         string
		err                 error
		path                string
		currentRelativePath = "."
	)
	filename := os.Getenv("env")
	if isTest {
		currentRelativePath = "../"
	}
	if currentPath, err = filepath.Abs(currentRelativePath); err != nil {
		panic(err)
	}
	path = filepath.Join(currentPath, "config")
	if filename == "" {
		path = filepath.Join(path, filePrefix+"dev.yml")
	} else {
		path = filepath.Join(path, filePrefix+filename+".yml")
	}
	if f, err := os.Open(path); err != nil {
		log.Fatalf("打开配置文件失败: %v", err)
	} else {
		if err := yaml.NewDecoder(f).Decode(config); err != nil {
			log.Fatalf("反序列化配置文件失败: %v", err)
		}
	}
}
