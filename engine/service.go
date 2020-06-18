package engine

import (
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-center/global"
	"github.com/dxvgef/tsing-center/load_balance"
)

// 设置本地数据中的服务
func SetService(service global.ServiceType) (err error) {
	if service.ID == "" {
		return errors.New("ID参数不能为空")
	}
	if service.LoadBalance == "" {
		return errors.New("LoadBalance参数不能为空")
	}

	mapValue, exist := global.Services.Load(service.ID)
	global.Services.Store(service.ID, service)
	if !exist {
		return nil
	}
	old, ok := mapValue.(global.ServiceType)
	if !ok {
		err = errors.New("类型断言失败")
		log.Error().Caller().Msg(err.Error())
		return nil
	}
	if old.LoadBalance == service.LoadBalance {
		return nil
	}
	// 缓存节点数据
	mapValue, exist = global.Nodes.Load(service.ID)
	if !exist {
		return nil
	}
	lb, lbOK := mapValue.(global.LoadBalance)
	if !lbOK {
		err = errors.New("类型断言失败")
		log.Error().Caller().Msg(err.Error())
		return nil
	}
	nodes := lb.Nodes()
	newNodes := make([]global.NodeType, len(nodes))
	for k := range nodes {
		newNodes[k].IP = nodes[k].IP
		newNodes[k].Port = nodes[k].Port
		newNodes[k].Weight = nodes[k].Weight
	}
	// 构建新的负载均衡实例
	newLB, err := load_balance.Build(service.LoadBalance)
	if err != nil {
		log.Err(err).Caller().Msg("构建负载均衡实例出错")
		return err
	}
	for k := range newNodes {
		newLB.Set(newNodes[k].IP, newNodes[k].Port, newNodes[k].Weight)
	}
	// 替换旧的实例
	global.Nodes.Store(service.ID, newLB)
	return nil
}

// 删除本地数据中的服务
func DelService(serviceID string) error {
	global.Services.Delete(serviceID)
	return nil
}

// 从本地数据中匹配服务
func MatchService(serviceID string) (global.ServiceType, bool) {
	if serviceID == "" {
		return global.ServiceType{}, false
	}
	mapValue, exist := global.Services.Load(serviceID)
	if !exist {
		return global.ServiceType{}, false
	}
	return mapValue.(global.ServiceType), true
}
