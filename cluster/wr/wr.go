package wr

import (
	"math/rand"
	"time"

	"github.com/dxvgef/tsing-center/global"
)

// 加权随机算法的集群
type Cluster struct {
	config      global.ServiceConfig
	nodes       []Node
	total       int
	totalWeight int        // 节点权重总和
	rand        *rand.Rand // 随机种子
}

type Node struct {
	ip      string // ip
	port    uint16 // 端口
	expires int64  // 生命周期截止时间(unix时间戳)
	weight  int    // 权重值
}

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
// 当weight的值<0>，表示不更新该属性
func (self *Cluster) Set(ip string, port uint16, weight int, expires int64) {
	for k := range self.nodes {
		// 如果节点已存在，则直接更新
		if self.nodes[k].ip == ip && self.nodes[k].port == port {
			if expires >= 0 && self.nodes[k].expires != expires {
				self.nodes[k].expires = expires
			}
			if self.nodes[k].weight >= 0 && self.nodes[k].weight != weight {
				self.nodes[k].weight = weight
				self.updateTotalWeight(weight)
				self.resetRand()
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
	self.updateTotalWeight(weight)
	self.resetRand()
	// 递增节点总数
	self.total++
}

// 获取节点总数
func (self *Cluster) Total() int {
	return self.total
}

// 移除节点
func (self *Cluster) Remove(ip string, port uint16) {
	for k := range self.nodes {
		if self.nodes[k].ip != ip && self.nodes[k].port != port {
			continue
		}
		self.nodes = append(self.nodes[:k], self.nodes[k+1:]...)
		self.updateTotalWeight(-self.nodes[k].weight)
		self.resetRand()
		self.total--
		return
	}
}

// 选举出下一个命中的节点
func (self *Cluster) Select() (string, uint16, int64) {
	if self.total == 0 {
		return "", 0, 0
	}
	if self.total == 1 {
		return self.nodes[0].ip, self.nodes[0].port, self.nodes[0].expires
	}
	if self.rand == nil {
		self.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	randomWeight := self.rand.Intn(self.totalWeight)
	for k := range self.nodes {
		randomWeight = randomWeight - self.nodes[k].weight
		if randomWeight <= 0 {
			return self.nodes[k].ip, self.nodes[k].port, self.nodes[k].expires
		}
	}
	return "", 0, 0
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

// 更新权重值总数
func (self *Cluster) updateTotalWeight(count int) {
	self.totalWeight += count
}

// 重设随机种子
func (self *Cluster) resetRand() {
	self.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
}
