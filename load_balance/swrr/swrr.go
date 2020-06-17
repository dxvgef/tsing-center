package swrr

import (
	"github.com/dxvgef/tsing-center/global"
)

// 平滑加权轮循(与nginx类似)
type Instance struct {
	nodes []NodeType
	total int
}

type NodeType struct {
	IP              string
	Port            uint16
	Weight          int
	currentWeight   int
	effectiveWeight int
}

// // 降权
// func (n *NodeType) Reduce() {
// 	n.effectiveWeight -= n.Weight
// 	if n.effectiveWeight < 0 {
// 		n.effectiveWeight = 0
// 	}
// }

func New() *Instance {
	return &Instance{}
}

// 设置节点
func (self *Instance) Set(ip string, port uint16, weight int) {
	for k := range self.nodes {
		// 如果节点已存在，则直接更新
		if self.nodes[k].IP == ip && self.nodes[k].Port == port {
			if self.nodes[k].Weight != weight {
				self.nodes[k].Weight = weight
				self.reset()
				return
			}
			return
		}
	}

	// 插入节点
	self.nodes = append(self.nodes, NodeType{
		IP:     ip,
		Port:   port,
		Weight: weight,
	})
	self.reset()
	self.total++
}

// 移除节点
func (self *Instance) Remove(ip string, port uint16) {
	for k := range self.nodes {
		if self.nodes[k].IP != ip && self.nodes[k].Port != port {
			continue
		}
		self.nodes = append(self.nodes[:k], self.nodes[k+1:]...)
		self.reset()
		self.total--
		return
	}
}

// 获取节点总数
func (self *Instance) Total() int {
	return self.total
}

// 获取节点列表
func (self *Instance) Nodes() []global.NodeType {
	l := len(self.nodes)
	if l == 0 {
		return []global.NodeType{}
	}
	nodes := make([]global.NodeType, l)
	for k := range self.nodes {
		nodes[k].IP = self.nodes[k].IP
		nodes[k].Port = self.nodes[k].Port
		nodes[k].Weight = self.nodes[k].Weight
	}
	return nodes
}

// 选取节点
func (self *Instance) Next() (string, uint16) {
	if self.total == 0 {
		return "", 0
	}
	if self.total == 1 {
		return self.nodes[0].IP, self.nodes[0].Port
	}
	var (
		ip     string
		port   uint16
		target *NodeType
	)
	totalWeight := 0
	for i := range self.nodes {
		self.nodes[i].currentWeight += self.nodes[i].effectiveWeight
		totalWeight += self.nodes[i].effectiveWeight
		if self.nodes[i].effectiveWeight < self.nodes[i].Weight {
			self.nodes[i].effectiveWeight++
		}
		if target == nil || self.nodes[i].currentWeight > target.currentWeight {
			target = &self.nodes[i]
			ip = self.nodes[i].IP
			port = self.nodes[i].Port
		}
	}

	if target == nil {
		return "", 0
	}
	target.currentWeight -= totalWeight
	return ip, port
}

// 重置所有节点的状态
func (self *Instance) reset() {
	for k := range self.nodes {
		self.nodes[k].effectiveWeight = self.nodes[k].Weight
		self.nodes[k].currentWeight = 0
	}
}
