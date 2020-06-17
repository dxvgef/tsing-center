package global

// 用于构建上游、中间件、存储器模块时的参数配置
type ModuleConfig struct {
	Name   string `json:"name"`
	Config string `json:"config"`
}

// 负载均衡接口
type LoadBalance interface {
	Set(string, uint16, int)
	Remove(string, uint16)
	Next() (string, uint16)
	Total() int
}

// 服务
type ServiceType struct {
	ID          string `json:"id"`                     // 服务ID
	LoadBalance string `json:"load_balance,omitempty"` // 负载均衡算法名称
}

// 节点
type NodeType struct {
	ServiceID string `json:"service_id"` // 服务ID
	IP        string `json:"ip"`         // 节点IP
	Port      uint16 `json:"port"`       // 节点端口
	Weight    int    `json:"weight"`     // 节点权重
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

	Watch() error // 监听存储器的数据变更
}
