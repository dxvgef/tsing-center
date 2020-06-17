package global

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	SnowflakeNode *snowflake.Node

	Storage StorageType // 存储器实例

	Services sync.Map // 所有服务列表,key=serviceID, value=loadBalance
	Nodes    sync.Map // 所有节点列表,key=serviceID, value=load_balance.Build()
)
