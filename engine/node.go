package engine

import (
	"errors"

	"github.com/dxvgef/tsing-center/global"
)

// 设置本地数据中的服务
func SetNode(serviceID string, node global.NodeType) error {
	if serviceID == "" {
		return errors.New("serviceID参数不能为空")
	}
	if node.IP == "" {
		return errors.New("node.IP参数不能为空")
	}
	if node.Port == 0 {
		return errors.New("node.Port参数不能为空")
	}
	service, exist := matchService(serviceID)
	if !exist {
		return errors.New("服务不存在")
	}

	global.Nodes.Store(serviceID, service)
	return nil
}

// 删除本地数据中的服务
func DelService(serviceID string) error {
	global.Services.Delete(serviceID)
	return nil
}

// 从本地数据中匹配服务
func matchService(serviceID string) (global.ServiceType, bool) {
	if serviceID == "" {
		return global.ServiceType{}, false
	}
	mapValue, exist := global.Services.Load(serviceID)
	if !exist {
		return global.ServiceType{}, false
	}
	return mapValue.(global.ServiceType), true
}
