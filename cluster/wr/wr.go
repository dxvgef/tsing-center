package wr

import (
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"

	"local/global"
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
	ttl     uint
	expires int64 // 生命周期截止时间(unix时间戳)
	weight  int   // 权重值
	meta    string
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
			node.TTL = self.nodes[k].ttl
			node.Expires = self.nodes[k].expires
			node.Weight = self.nodes[k].weight
			node.Mete = self.nodes[k].meta
			return
		}
	}
	return
}

// 设置节点
// 当weight的值<0>，表示不更新该属性
func (self *Cluster) Set(node global.Node) {
	for k := range self.nodes {
		// 如果节点已存在，则直接更新
		if self.nodes[k].ip == node.IP && self.nodes[k].port == node.Port {
			self.nodes[k].ttl = node.TTL
			self.nodes[k].expires = node.Expires
			if self.nodes[k].weight >= 0 && self.nodes[k].weight != node.Weight {
				self.nodes[k].weight = node.Weight
				self.updateTotalWeight(node.Weight)
				self.resetRand()
			}
			self.nodes[k].meta = node.Mete
			return
		}
	}

	// 插入节点
	self.nodes = append(self.nodes, Node{
		ip:      node.IP,
		port:    node.Port,
		weight:  node.Weight,
		ttl:     node.TTL,
		expires: node.Expires,
		meta:    node.Mete,
	})
	self.updateTotalWeight(node.Weight)
	self.resetRand()
	// 递增节点总数
	self.total++
}

// 触活节点
func (self *Cluster) Touch(ip string, port uint16, expires int64) {
	for k := range self.nodes {
		if self.nodes[k].ip == ip && self.nodes[k].port == port && self.nodes[k].ttl > 0 {
			self.nodes[k].expires = expires
			return
		}
	}
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
func (self *Cluster) Select() (node global.Node) {
	var lostNodes []global.Node

	defer func() {
		if len(lostNodes) > 0 {
			go func() {
				if err := global.Storage.Clean(self.config.ServiceID, lostNodes); err != nil {
					log.Err(err).Caller().Send()
					return
				}
			}()
		}
	}()

	now := time.Now().Unix()

	switch self.total {
	case 0:
		return
	case 1:
		if self.nodes[0].weight < 0 || (self.nodes[0].expires > 0 && self.nodes[0].expires <= now) {
			lostNodes = append(lostNodes, global.Node{
				IP:   self.nodes[0].ip,
				Port: self.nodes[0].port,
			})
			return
		}
		node.IP = self.nodes[0].ip
		node.Port = self.nodes[0].port
		node.TTL = self.nodes[0].ttl
		node.Expires = self.nodes[0].expires
		node.Mete = self.nodes[0].meta
		return
	}
	if self.rand == nil {
		self.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	var randomWeight = self.rand.Intn(self.totalWeight)

	for i := range self.nodes {
		if self.nodes[i].weight == 0 {
			lostNodes = append(lostNodes, global.Node{
				IP:   self.nodes[i].ip,
				Port: self.nodes[i].port,
			})
			continue
		}
		if self.nodes[i].expires > 0 && self.nodes[i].expires <= now {
			lostNodes = append(lostNodes, global.Node{
				IP:   self.nodes[i].ip,
				Port: self.nodes[i].port,
			})
			continue
		}
		randomWeight = randomWeight - self.nodes[i].weight
		if randomWeight <= 0 {
			node.IP = self.nodes[i].ip
			node.Port = self.nodes[i].port
			node.TTL = self.nodes[i].ttl
			node.Expires = self.nodes[i].expires
			node.Mete = self.nodes[i].meta
			break
		}
	}

	return
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
		nodes[k].TTL = self.nodes[k].ttl
		nodes[k].Weight = self.nodes[k].weight
		nodes[k].Expires = self.nodes[k].expires
		nodes[k].Mete = self.nodes[k].meta
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
