package global

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	SnowflakeNode *snowflake.Node

	Storage StorageType // 存储器实例

	// 服务列表，服务是个集群的抽象概念
	// key=服务ID, value=集群实例
	Services      sync.Map
	TotalServices uint32 // 服务总数
)

// 服务配置，用作集群构建时的参数
type ServiceConfig struct {
	ServiceID   string `json:"service_id"`   // 服务ID
	LoadBalance string `json:"load_balance"` // 负载
}

// 节点属性
type Node struct {
	IP      string `json:"ip"`      // 节点IP
	Port    uint16 `json:"port"`    // 节点端口
	Weight  int    `json:"weight"`  // 节点权重
	TTL     uint   `json:"ttl"`     // TTL(秒)
	Expires int64  `json:"expires"` // 生命周期截止时间(unix时间戳)，值为0表示一直有效
}

// 集群接口
type Cluster interface {
	Config() ServiceConfig       // 获得配置
	Set(Node)                    // 设置节点
	Touch(string, uint16, int64) // 触活节点，入参(ip, port, expires)
	Remove(string, uint16)       // 移除节点，入参(ip, port)
	Select() Node                // 选取节点
	Total() int                  // 节点总数
	Nodes() []Node               // 节点列表
	Find(string, uint16) Node    // 查找某个节点，入参(ip, port)
}

// 存储器接口
type StorageType interface {
	LoadAll() error // 从存储器加载所有数据到本地
	SaveAll() error // 将本地所有数据保存到存储器

	LoadService([]byte) error          // 从存储器加载单个服务数据，入参(ServiceConfig的json字节码)
	SaveService(ServiceConfig) error   // 将本地单个服务保存到存储器
	DeleteLocalService(string) error   // 删除本地单个服务，入参(存储器key)
	DeleteStorageService(string) error // 删除存储器中单个服务

	LoadNode(string, []byte) error                  // 从存储器加载单个节点数据，入参(存储器key，存储器数据)
	SaveNode(string, Node) error                    // 将本地单个节点保存到存储器
	DeleteLocalNode(string) error                   // 删除本地单个节点，入参(存储器key)
	DeleteStorageNode(string, string, uint16) error // 删除存储器中单个节点，入参(服务id, ip, port)

	Clean(string, []Node) error // 清理已失效的节点

	Watch() // 监听存储器的数据变更
}
