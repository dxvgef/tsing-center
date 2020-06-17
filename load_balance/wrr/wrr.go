package wrr

// 加权轮循(与LVS类似)
type Instance struct {
	nodes         []NodeType // 节点列表
	total         int        // 节点总数
	weightGCD     int        // 权总值最大公约数
	maxWeight     int        // 最大权重值
	lastIndex     int        // 最后命中的节点索引，初始值是-1
	currentWeight int        // 当前权重值
}

type NodeType struct {
	IP     string // 地址
	Port   uint16 // 端口
	Weight int    // 权重值
}

func New() *Instance {
	return &Instance{
		lastIndex: -1,
	}
}

// 设置节点
func (self *Instance) Set(ip string, port uint16, weight int) {
	for k := range self.nodes {
		if self.nodes[k].IP == ip && self.nodes[k].Port == port {
			self.nodes[k].Weight = weight
			self.calcMaxWeight(weight)
			return
		}
	}
	self.nodes = append(self.nodes, NodeType{
		IP:     ip,
		Port:   port,
		Weight: weight,
	})
	self.total++
	self.calcMaxWeight(weight)
}

// 获得节点总数
func (self *Instance) Total() int {
	return self.total
}

// 移除节点
func (self *Instance) Remove(ip string, port uint16) {
	for k := range self.nodes {
		if self.nodes[k].IP != ip && self.nodes[k].Port != port {
			continue
		}
		self.calcAllGCD(self.total)
		self.total--
		self.calcMaxWeight(self.nodes[k].Weight)
		self.nodes = append(self.nodes[:k], self.nodes[k+1:]...)
		break
	}
}

// 选举节点
func (self *Instance) Next() (string, uint16) {
	switch self.total {
	case 0:
		return "", 0
	case 1:
		return self.nodes[0].IP, self.nodes[0].Port
	}

	var cw, gcd, last, total int
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
					return "", 0
				}
			}
		}
		cw = self.currentWeight
		if self.nodes[last].Weight >= cw {
			return self.nodes[last].IP, self.nodes[last].Port
		}
	}
}

// 计算最大权重值
func (self *Instance) calcMaxWeight(weight int) {
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
func (self *Instance) calcAllGCD(i int) int {
	if i == 1 {
		return self.nodes[0].Weight
	}
	return calcGCD(self.nodes[i-1].Weight, self.calcAllGCD(i-1))
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
