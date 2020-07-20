package engine

import (
	"errors"

	"local/global"
)

// 设置本地数据中的节点
func SetNode(serviceID string, node global.Node) (err error) {
	if serviceID == "" {
		return errors.New("serviceID参数不能为空")
	}
	if node.IP == "" {
		return errors.New("ip参数不能为空")
	}
	if node.Port == 0 {
		return errors.New("port参数不能为0")
	}

	// 获取集群实例
	ci := FindCluster(serviceID)
	if ci == nil {
		return errors.New("服务不存在或不可用")
	}
	ci.Set(node)
	return nil
}

// 删除本地数据中的节点
func DelNode(serviceID string, ip string, port uint16) error {
	ci := FindCluster(serviceID)
	if ci == nil {
		return errors.New("服务不存在或不可用")
	}
	ci.Remove(ip, port)
	if ci.Total() == 0 {
		global.Services.Delete(serviceID)
		return nil
	}
	return nil
}
