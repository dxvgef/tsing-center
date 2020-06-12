package engine

import (
	"errors"

	"github.com/dxvgef/tsing-center/global"
)

// 设置本地数据中的服务
func SetService(service global.ServiceType) error {
	if service.ID == "" {
		return errors.New("ID参数不能为空")
	}
	if service.LoadBalance == "" {
		return errors.New("LoadBalance参数不能为空")
	}
	global.Services.Store(service.ID, service)
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
