package global

// 用于构建上游、中间件、存储器模块时的参数配置
type ModuleConfig struct {
	Name   string `json:"name"`   // 模块名称
	Config string `json:"config"` // 配置(json字符串)
}

// 负载均衡接口
type LoadBalance interface {
	Name() string            // 算法名称
	Set(string, uint16, int) // 设置节点
	Remove(string, uint16)   // 移除节点
	Next() (string, uint16)  // 选取节点
	Total() int              // 节点总数
	Nodes() []NodeType       // 节点列表
}

// 服务
type ServiceType struct {
	ID          string `json:"id"`           // 服务ID
	LoadBalance string `json:"load_balance"` // 负载均衡算法名称
}

// 节点
type NodeType struct {
	IP     string `json:"ip"`     // 节点IP
	Port   uint16 `json:"port"`   // 节点端口
	Weight int    `json:"weight"` // 节点权重
}

// 存储器
type StorageType interface {
	LoadAll() error // 从存储器加载所有数据到本地
	SaveAll() error // 将本地所有数据保存到存储器

	LoadAllService() error             // 从存储器加载所有服务到本地
	LoadService(string, []byte) error  // 从存储器加载单个服务数据
	SaveAllService() error             // 将本地所有服务保存到存储器
	SaveService(string, string) error  // 将本地单个服务保存到存储器
	DeleteLocalService(string) error   // 删除本地单个服务
	DeleteStorageService(string) error // 删除存储器中单个服务

	LoadAllNode() error                             // 从存储器加载所有节点到本地
	LoadNode(string, []byte) error                  // 从存储器加载单个节点数据
	SaveAllNode() error                             // 将本地所有节点保存到存储器
	SaveNode(string, string, uint16, int) error     // 将本地单个节点保存到存储器
	DeleteLocalNode(string) error                   // 删除本地单个节点
	DeleteStorageNode(string, string, uint16) error // 删除存储器中单个节点

	Watch() error // 监听存储器的数据变更
}
