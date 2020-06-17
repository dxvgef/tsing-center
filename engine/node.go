package engine

import (
	"errors"

	"github.com/dxvgef/tsing-center/global"
	"github.com/dxvgef/tsing-center/load_balance"
)

// 设置本地数据中的节点
func SetNode(serviceID string, ip string, port uint16, weight int) (err error) {
	if serviceID == "" {
		return errors.New("serviceID参数不能为空")
	}
	if ip == "" {
		return errors.New("ip参数不能为空")
	}
	if port == 0 {
		return errors.New("port参数不能为0")
	}
	service, serviceExist := MatchService(serviceID)
	if !serviceExist {
		return errors.New("服务ID不存在")
	}
	lb, lbExist := matchNode(serviceID)
	if !lbExist {
		lb, err = load_balance.Build(service.LoadBalance)
		if err != nil {
			return err
		}
	}
	lb.Set(ip, port, weight)
	global.Nodes.Store(serviceID, lb)
	return nil
}

// 删除本地数据中的节点
func DelNode(serviceID string, ip string, port uint16) error {
	mapValue, exist := global.Nodes.Load(serviceID)
	if !exist {
		return nil
	}
	lb, ok := mapValue.(global.LoadBalance)
	if !ok {
		return errors.New("未知的负载均衡类型")
	}
	lb.Remove(ip, port)
	if lb.Total() == 0 {
		global.Nodes.Delete(serviceID)
		return nil
	}
	global.Nodes.Store(serviceID, lb)
	return nil
}

// 判断节点是否存在
func NodeExist(serviceID, ip string, port uint16) bool {
	mapValue, exist := global.Nodes.Load(serviceID)
	if !exist {
		return false
	}
	lb, ok := mapValue.(global.LoadBalance)
	if !ok {
		return false
	}
	nodes := lb.Nodes()
	for k := range nodes {
		if nodes[k].IP == ip && nodes[k].Port == port {
			return true
		}
	}
	return false
}

// 从本地数据中匹配节点
func matchNode(serviceID string) (global.LoadBalance, bool) {
	mapValue, exist := global.Nodes.Load(serviceID)
	if !exist {
		return nil, false
	}
	lb, ok := mapValue.(global.LoadBalance)
	return lb, ok
}
