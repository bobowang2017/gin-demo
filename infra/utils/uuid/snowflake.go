package uuid

import (
	"gin-demo/infra/utils/log"
	"github.com/bwmarrin/snowflake"
	"sync"
)

var (
	snowflakeNodeClient     *snowflake.Node
	snowflakeNodeClientOnce sync.Once
)

func SnowFlakeSetUp(serverNodeId int64) {
	snowflake.Epoch = 1609430401000
	node, err := snowflake.NewNode(serverNodeId)
	if err != nil {
		log.Logger.Errorf("snowflake NewNode error: %s", err)
		return
	}
	snowflakeNodeClient = node
	return
}

func NewSnowflakeClient() *snowflake.Node {
	snowflakeNodeClientOnce.Do(func() {
		SnowFlakeSetUp(0)
	})
	return snowflakeNodeClient
}

func SnowFlakeId() int64 {
	return NewSnowflakeClient().Generate().Int64()
}
