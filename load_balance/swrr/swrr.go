package swrr

import (
	"errors"
	"sync"

	"github.com/rs/zerolog/log"
)

var Inst *InstType
var globalNodes sync.Map // 节点列表 key=upstreamID, value=[]*NodeType
var globalTotal sync.Map // 节点总量 key=upstreamID, value=int

type InstType struct{}

// 节点结构
type NodeType struct {
	IP              string
	Port            uint32
	Weight          int
	CurrentWeight   int
	EffectiveWeight int
}

func Init() *InstType {
	if Inst != nil {
		return Inst
	}
	return &InstType{}
}

// // 降权
// func (n *NodeType) Reduce() {
// 	n.EffectiveWeight -= n.Weight
// 	if n.EffectiveWeight < 0 {
// 		n.EffectiveWeight = 0
// 	}
// }

// 设置单个节点
func (self *InstType) Set(serviceID, ip string, port uint32, weight int) (err error) {
	var nodes []*NodeType
	mapValue, exist := globalNodes.Load(serviceID)
	if exist {
		var ok bool
		if nodes, ok = mapValue.([]*NodeType); !ok {
			err = errors.New("类型断言失败")
			log.Err(err).Caller().Msg("设置节点")
			return
		}
	}
	for k := range nodes {
		// 如果节点已存在，则直接更新
		if nodes[k].IP == ip && nodes[k].Port == port {
			if nodes[k].Weight != weight {
				nodes[k].Weight = weight
				globalNodes.Store(serviceID, nodes)
				return self.reset(serviceID)
			}
			return
		}
	}

	// 插入节点
	nodes = append(nodes, &NodeType{
		IP:     ip,
		Port:   port,
		Weight: weight,
	})
	globalNodes.Store(serviceID, nodes)
	if err = self.reset(serviceID); err != nil {
		return
	}
	// 递增节点总数
	return self.updateTotal(serviceID, 1)
}

// 移除单个节点
func (self *InstType) Remove(serviceID, ip string, port uint32) (err error) {
	mapValue, exist := globalNodes.Load(serviceID)
	if !exist {
		return nil
	}
	if ip == "" {

		globalNodes.Delete(serviceID)
	}
	nodes, ok := mapValue.([]*NodeType)
	if !ok {
		err = errors.New("类型断言失败")
		log.Err(err).Caller().Msg("移除节点")
		return
	}
	for k := range nodes {
		if nodes[k].IP != ip && nodes[k].Port != port {
			continue
		}
		newNodes := append(nodes[:k], nodes[k+1:]...)
		globalNodes.Store(serviceID, newNodes)
		if err = self.reset(serviceID); err != nil {
			return
		}
		return self.updateTotal(serviceID, -1)
	}
	return
}

// 获取节点总数
func (self *InstType) Total(serviceID string) int {
	mapValue, exist := globalTotal.Load(serviceID)
	if !exist {
		return 0
	}
	total, ok := mapValue.(int)
	if !ok {
		return 0
	}
	return total
}

// 选取节点
func (self *InstType) Next(serviceID string) (string, uint32) {
	nodes := self.getNodes(serviceID)
	if nodes == nil {
		return "", 0
	}
	nodeTotal := len(nodes)
	if nodeTotal == 1 {
		return nodes[0].IP, nodes[0].Port
	}
	var (
		ip     string
		port   uint32
		target *NodeType
	)
	totalWeight := 0
	for i := range nodes {
		nodes[i].CurrentWeight += nodes[i].EffectiveWeight
		totalWeight += nodes[i].EffectiveWeight
		if nodes[i].EffectiveWeight < nodes[i].Weight {
			nodes[i].EffectiveWeight++
		}
		if target == nil || nodes[i].CurrentWeight > target.CurrentWeight {
			target = nodes[i]
			ip = nodes[i].IP
			port = nodes[i].Port
		}
	}

	globalNodes.Store(serviceID, nodes)

	if target == nil {
		return "", 0
	}
	target.CurrentWeight -= totalWeight
	return ip, port
}

// 重置所有节点的状态
func (self *InstType) reset(serviceID string) (err error) {
	mapValue, exist := globalNodes.Load(serviceID)
	if !exist {
		return nil
	}
	nodes, ok := mapValue.([]*NodeType)
	if !ok {
		err = errors.New("类型断言失败")
		log.Err(err).Caller().Msg("重置节点")
		return
	}
	for k := range nodes {
		nodes[k].EffectiveWeight = nodes[k].Weight
		nodes[k].CurrentWeight = 0
	}
	return nil
}

// 更新节点总数
func (self *InstType) updateTotal(serviceID string, count int) (err error) {
	mapValue, exist := globalTotal.Load(serviceID)
	if !exist {
		globalTotal.Store(serviceID, 0)
		return nil
	}

	total, ok := mapValue.(int)
	if !ok {
		err = errors.New("类型断言失败")
		log.Err(err).Caller().Msg("更新节点总数")
		return err
	}
	total += count
	globalTotal.Store(serviceID, total)
	return nil
}

// 判断节点是否存在
func (self *InstType) nodeExist(serviceID, ip string, port uint32) (exist bool) {
	if _, exist = globalNodes.Load(serviceID); !exist {
		return
	}
	globalNodes.Range(func(key, value interface{}) bool {
		if key.(string) == serviceID {
			nodes, ok := value.([]*NodeType)
			if !ok {
				return true
			}
			for k := range nodes {
				if nodes[k].IP == ip && nodes[k].Port == port {
					exist = true
					return false
				}
			}
		}
		return true
	})
	return
}

// 获得所有节点
func (self *InstType) getNodes(serviceID string) []*NodeType {
	mapValue, exist := globalNodes.Load(serviceID)
	if !exist {
		return nil
	}
	nodes, ok := mapValue.([]*NodeType)
	if !ok {
		return nil
	}
	return nodes
}
