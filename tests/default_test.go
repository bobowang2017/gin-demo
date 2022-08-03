package tests

import (
	"gin-demo/core/settings"
	"gin-demo/infra/model"
	"gin-demo/infra/utils/log"
	"testing"
)

func SetUp() {
	settings.SetupTest()
	log.Setup()
	model.Setup()
}

func TestDB(t *testing.T) {
	SetUp()
}
