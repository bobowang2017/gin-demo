package main

import (
	"flag"
	"gin-demo/core"
	"gin-demo/infra/common"
	"gin-demo/infra/utils/log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	modulesStr := flag.String("m", "core", "模块启动选项，默认core启动")
	env := flag.String("e", "local", "配置环境，默认dev，该参数将影响挂载配置文件名")
	flag.Parse()

	if err := os.Setenv("ENV", *env); err != nil {
		log.Logger.Warnf("设置启动配置环境失败，默认使用dev配置环境 %v", err)
	}

	modules := make(map[string]string)

	// 去重
	for _, item := range strings.Split(*modulesStr, ",") {
		module := strings.Trim(item, " ")
		modules[module] = module
	}

	for module := range modules {
		switch strings.Trim(module, " ") {
		case common.ModuleCore:
			go core.Start()
		}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, os.Interrupt, syscall.SIGTERM)
	// Wait for OS termination signal
	<-sigChan
}
