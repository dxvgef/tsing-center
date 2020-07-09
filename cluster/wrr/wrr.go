package wrr

import (
	"time"

	"github.com/rs/zerolog/log"

	"local/global"
)

// 加权轮循(与LVS类似)算法的集群
type Cluster struct {
	config        global.ServiceConfig
	nodes         []Node // 节点列表
	total         int    // 节点总数
	weightGCD     int    // 权总值最大公约数
	maxWeight     int    // 最大权重值
	lastIndex     int    // 最后命中的节点索引，初始值是-1
	currentWeight int    // 当前权重值
}

type Node struct {
	ip      string // 地址
	port    uint16 // 端口
	expires int64  // 生命周期截止时间(unix时间戳)
	weight  int    // 权重值
}

func New(config global.ServiceConfig) *Cluster {
	return &Cluster{
		config:    config,
		lastIndex: -1,
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
		if self.nodes[k].ip == ip && self.nodes[k].port == port {
			if expires >= 0 && self.nodes[k].expires != expires {
				self.nodes[k].expires = expires
			}
			if self.nodes[k].weight >= 0 && self.nodes[k].weight != weight {
				self.nodes[k].weight = weight
				self.calcMaxWeight(weight)
				self.calcAllGCD(self.total)
			}
			return
		}
	}
	self.nodes = append(self.nodes, Node{
		ip:      ip,
		port:    port,
		weight:  weight,
		expires: expires,
	})
	self.total++
	self.calcMaxWeight(weight)
	self.calcAllGCD(self.total)
}

// 获得节点总数
func (self *Cluster) Total() int {
	return self.total
}

// 移除节点
func (self *Cluster) Remove(ip string, port uint16) {
	for k := range self.nodes {
		if self.nodes[k].ip != ip && self.nodes[k].port != port {
			continue
		}
		self.calcAllGCD(self.total)
		self.total--
		self.calcMaxWeight(self.nodes[k].weight)
		self.nodes = append(self.nodes[:k], self.nodes[k+1:]...)
		break
	}
}

// 选举节点
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

	var (
		cw, gcd, last, total int
		lostNodes            []global.Node
		now                  = time.Now().Unix()
	)
	for {
		last = self.lastIndex
		total = self.total
		last = (last + 1) % total
		self.lastIndex = last
		if last == 0 {
			cw = self.currentWeight
			gcd = self.weightGCD
			newCW := cw - gcd
			self.currentWeight = newCW
			if newCW <= 0 {
				newCW = self.maxWeight
				self.currentWeight = newCW
				if newCW == 0 {
					break
				}
			}
		}
		if self.nodes[last].weight == 0 {
			continue
		}
		if self.nodes[last].expires > 0 && self.nodes[last].expires <= now {
			lostNodes = append(lostNodes, global.Node{
				IP:   self.nodes[last].ip,
				Port: self.nodes[last].port,
			})
			continue
		}
		cw = self.currentWeight
		if self.nodes[last].weight >= cw {
			ip = self.nodes[last].ip
			port = self.nodes[last].port
			expires = self.nodes[last].expires
			break
		}
	}
	if len(lostNodes) > 0 {
		go func() {
			if err := global.Storage.Clean(self.config.ServiceID, lostNodes); err != nil {
				log.Err(err).Caller().Send()
				return
			}
		}()
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
		nodes[k].Weight = self.nodes[k].weight
		nodes[k].Expires = self.nodes[k].expires
	}
	return nodes
}

// 计算最大权重值
func (self *Cluster) calcMaxWeight(weight int) {
	if weight == 0 {
		return
	}
	if self.weightGCD == 0 {
		self.weightGCD = weight
		self.maxWeight = weight
		self.lastIndex = -1
		self.currentWeight = 0
		return
	}
	self.weightGCD = calcGCD(self.weightGCD, weight)
	if self.maxWeight < weight {
		self.maxWeight = weight
	}
}

// 计算所有节点权重值的最大公约数
func (self *Cluster) calcAllGCD(i int) int {
	if i == 1 {
		return self.nodes[0].weight
	}
	return calcGCD(self.nodes[i-1].weight, self.calcAllGCD(i-1))
}

// 计算两个权重值的最大公约数
func calcGCD(a, b int) int {
	if a < b {
		a, b = b, a // 交换a和b
	}
	if b == 0 {
		return a
	}
	return calcGCD(b, a%b)
}
