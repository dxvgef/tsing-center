package swrr

import (
	"time"

	"github.com/dxvgef/tsing-center/global"
)

// 平滑加权轮循(与nginx类似)算法的集群
type Cluster struct {
	config global.ServiceConfig
	nodes  []Node
	total  int
}

type Node struct {
	ip              string
	port            uint16
	expires         int64 // 生命周期截止时间(unix时间戳)
	weight          int
	currentWeight   int
	effectiveWeight int
}

// // 降权
// func (n *Node) Reduce() {
// 	n.effectiveWeight -= n.weight
// 	if n.effectiveWeight < 0 {
// 		n.effectiveWeight = 0
// 	}
// }

func New(config global.ServiceConfig) *Cluster {
	return &Cluster{
		config: config,
	}
}

// 获得配置
func (self *Cluster) Config() global.ServiceConfig {
	return self.config
}

// 查找某个节点
func (self *Cluster) Find(ip string, port uint16) (node global.Node) {
	for k := range self.nodes {
		if self.nodes[k].ip == ip && self.nodes[k].port == port {
			node.IP = ip
			node.Port = port
			node.Expires = self.nodes[k].expires
			node.Weight = self.nodes[k].weight
			return
		}
	}
	return
}

// 设置节点
// 当weight的值<0，表示不更新该属性
func (self *Cluster) Set(ip string, port uint16, weight int, expires int64) {
	for k := range self.nodes {
		// 如果节点已存在，则直接更新
		if self.nodes[k].ip == ip && self.nodes[k].port == port {
			if expires >= 0 && self.nodes[k].expires != expires {
				self.nodes[k].expires = expires
			}
			if self.nodes[k].weight >= 0 && self.nodes[k].weight != weight {
				self.nodes[k].weight = weight
				self.reset()
			}
			return
		}
	}

	// 插入节点
	self.nodes = append(self.nodes, Node{
		ip:      ip,
		port:    port,
		weight:  weight,
		expires: expires,
	})
	self.reset()
	self.total++
}

// 移除节点
func (self *Cluster) Remove(ip string, port uint16) {
	for k := range self.nodes {
		if self.nodes[k].ip != ip && self.nodes[k].port != port {
			continue
		}
		self.nodes = append(self.nodes[:k], self.nodes[k+1:]...)
		self.reset()
		self.total--
		return
	}
}

// 获取节点总数
func (self *Cluster) Total() int {
	return self.total
}

// 获取节点列表
func (self *Cluster) Nodes() []global.Node {
	l := len(self.nodes)
	if l == 0 {
		return []global.Node{}
	}
	nodes := make([]global.Node, l)
	for k := range self.nodes {
		nodes[k].IP = self.nodes[k].ip
		nodes[k].Port = self.nodes[k].port
		nodes[k].Weight = self.nodes[k].weight
		nodes[k].Expires = self.nodes[k].expires
	}
	return nodes
}

// 选取节点
func (self *Cluster) Select() (ip string, port uint16, expires int64) {
	switch self.total {
	case 0:
		return
	case 1:
		ip = self.nodes[0].ip
		port = self.nodes[0].port
		expires = self.nodes[0].expires
		return
	}
	var target *Node
	totalWeight := 0
	now := time.Now().Unix()
	var lostNodes []global.Node
	for i := range self.nodes {
		if self.nodes[i].expires <= now {
			lostNodes = append(lostNodes, global.Node{
				IP:   self.nodes[i].ip,
				Port: self.nodes[i].port,
			})
			continue
		}
		if self.nodes[i].weight == 0 {
			continue
		}
		self.nodes[i].currentWeight += self.nodes[i].effectiveWeight
		totalWeight += self.nodes[i].effectiveWeight
		if self.nodes[i].effectiveWeight < self.nodes[i].weight {
			self.nodes[i].effectiveWeight++
		}
		if target == nil || self.nodes[i].currentWeight > target.currentWeight {
			target = &self.nodes[i]
			ip = self.nodes[i].ip
			port = self.nodes[i].port
			expires = self.nodes[i].expires
		}
	}
	if target == nil {
		return
	}
	target.currentWeight -= totalWeight
	if len(lostNodes) > 0 {
		go global.Storage.Clean(self.config.ServiceID, lostNodes)
	}
	return
}

// 重置所有节点的状态
func (self *Cluster) reset() {
	for k := range self.nodes {
		self.nodes[k].effectiveWeight = self.nodes[k].weight
		self.nodes[k].currentWeight = 0
	}
}
