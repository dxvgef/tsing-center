package wr

import (
	"math/rand"
	"time"

	"github.com/dxvgef/tsing-center/global"
)

// 加权随机
type Instance struct {
	nodes       []NodeType
	total       int
	totalWeight int        // 节点权重总和
	rand        *rand.Rand // 随机种子
}

type NodeType struct {
	IP     string // IP
	Port   uint16 // 端口
	Weight int    // 权重值
}

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
				self.updateTotalWeight(weight)
				self.resetRand()
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
	self.updateTotalWeight(weight)
	self.resetRand()
	// 递增节点总数
	self.total++
}

// 获取节点总数
func (self *Instance) Total() int {
	return self.total
}

// 移除节点
func (self *Instance) Remove(ip string, port uint16) {
	for k := range self.nodes {
		if self.nodes[k].IP != ip && self.nodes[k].Port != port {
			continue
		}
		self.nodes = append(self.nodes[:k], self.nodes[k+1:]...)
		self.updateTotalWeight(-self.nodes[k].Weight)
		self.resetRand()
		self.total--
		return
	}
}

// 选举出下一个命中的节点
func (self *Instance) Next() (string, uint16) {
	if self.total == 0 {
		return "", 0
	}
	if self.total == 1 {
		return self.nodes[0].IP, self.nodes[0].Port
	}
	if self.rand == nil {
		self.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	randomWeight := self.rand.Intn(self.totalWeight)
	for k := range self.nodes {
		randomWeight = randomWeight - self.nodes[k].Weight
		if randomWeight <= 0 {
			return self.nodes[k].IP, self.nodes[k].Port
		}
	}
	return "", 0
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

// 更新权重值总数
func (self *Instance) updateTotalWeight(count int) {
	self.totalWeight += count
}

// 重设随机种子
func (self *Instance) resetRand() {
	self.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
}
