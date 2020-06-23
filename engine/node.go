package engine

import (
	"errors"

	"local/global"
)

// 设置本地数据中的节点
func SetNode(serviceID string, ip string, port uint16, weight int, expires int64) (err error) {
	if serviceID == "" {
		return errors.New("serviceID参数不能为空")
	}
	if ip == "" {
		return errors.New("ip参数不能为空")
	}
	if port == 0 {
		return errors.New("port参数不能为0")
	}

	// 获取集群实例
	ci := FindCluster(serviceID)
	if ci == nil {
		return errors.New("服务不存在或不可用")
	}
	ci.Set(ip, port, weight, expires)
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
