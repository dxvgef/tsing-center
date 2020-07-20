package engine

import (
	"errors"
	"sync/atomic"

	"github.com/rs/zerolog/log"

	"local/cluster"
	"local/global"
)

// 设置本地数据中的服务
func SetService(config global.ServiceConfig) (err error) {
	if config.ServiceID == "" {
		return errors.New("ID参数不能为空")
	}
	if config.LoadBalance == "" {
		return errors.New("LoadBalance参数不能为空")
	}
	var newCluster global.Cluster
	// 获取旧的集群实例
	oldCluster := FindCluster(config.ServiceID)
	// 如果原来的集群无效
	if oldCluster == nil {
		// 构建新的集群实例
		newCluster, err = cluster.Build(config)
		if err != nil {
			log.Err(err).Caller().Send()
			return
		}
		// 写入本地服务列表
		global.Services.Store(config.ServiceID, newCluster)
		addTotalServices(1)
		return nil
	}

	// 缓存节点数据
	nodes := oldCluster.Nodes()
	// 构建新的负载均衡实例
	newCluster, err = cluster.Build(config)
	if err != nil {
		log.Err(err).Caller().Send()
		return
	}
	// 将缓存中的节点写入到新的集群实例中
	for k := range nodes {
		newCluster.Set(global.Node{
			IP:      nodes[k].IP,
			Port:    nodes[k].Port,
			Weight:  nodes[k].Weight,
			TTL:     nodes[k].TTL,
			Expires: nodes[k].Expires,
		})
	}
	// 替换旧的集群实例
	global.Services.Store(config.ServiceID, newCluster)
	return nil
}

// 删除本地数据中的服务
func DelService(serviceID string) error {
	global.Services.Delete(serviceID)
	addTotalServices(-1)
	return nil
}

// 从本地数据中匹配集群实例
func FindCluster(serviceID string) (ci global.Cluster) {
	mapValue, exist := global.Services.Load(serviceID)
	if !exist {
		return nil
	}
	cluster, ok := mapValue.(global.Cluster)
	if !ok {
		return nil
	}
	return cluster
}

// 递值global.TotalServices的值，如果v是负数则递减
func addTotalServices(v int) {
	atomic.AddUint32(&global.TotalServices, ^uint32(-v-1))
}
