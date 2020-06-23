package etcd

import (
	"context"
	"errors"
	"path"
	"strings"
	"time"

	"github.com/dxvgef/tsing-center/engine"
	"github.com/dxvgef/tsing-center/global"

	"github.com/rs/zerolog/log"
)

// 从存储器加载服务到本地，如果不存在则创建
func (self *Etcd) LoadService(data []byte) (err error) {
	var service global.ServiceConfig
	if err = service.UnmarshalJSON(data); err != nil {
		log.Err(err).Caller().Send()
		return
	}
	return engine.SetService(service)
}

// 将本地服务数据保存到存储器中，如果不存在则创建
func (self *Etcd) SaveService(config global.ServiceConfig) (err error) {
	var configBytes []byte
	configBytes, err = config.MarshalJSON()
	if err != nil {
		log.Err(err).Caller().Send()
		return
	}

	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/services/")
	key.WriteString(global.EncodeKey(config.ServiceID))
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err = self.client.Put(ctx, key.String(), global.BytesToStr(configBytes)); err != nil {
		log.Err(err).Caller().Send()
		return err
	}
	return nil
}

// 删除本地服务数据
func (self *Etcd) DeleteLocalService(key string) error {
	serviceID, err := global.DecodeKey(path.Base(key))
	if err != nil {
		log.Err(err).Caller().Send()
		return err
	}
	return engine.DelService(serviceID)
}

// 删除存储器中服务数据
func (self *Etcd) DeleteStorageService(serviceID string) error {
	if serviceID == "" {
		return errors.New("服务ID不能为空")
	}

	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/services/")
	key.WriteString(serviceID)
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	_, err := self.client.Delete(ctx, key.String())
	if err != nil {
		log.Err(err).Caller().Send()
		return err
	}
	return nil
}
