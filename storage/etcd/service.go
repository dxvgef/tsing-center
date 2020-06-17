package etcd

import (
	"context"
	"errors"
	"path"
	"strings"
	"time"

	"github.com/dxvgef/tsing-center/engine"
	"github.com/dxvgef/tsing-center/global"

	"github.com/coreos/etcd/clientv3"
	"github.com/rs/zerolog/log"
)

// 从存储器加载服务到本地，如果不存在则创建
func (self *Etcd) LoadService(key string, data []byte) error {
	serviceID, err := global.DecodeKey(path.Base(key))
	if err != nil {
		return err
	}
	var service global.ServiceType
	service.ID = serviceID
	service.LoadBalance = global.BytesToStr(data)
	return engine.SetService(service)
}

// 从存储器加载所有主机数据到本地
func (self *Etcd) LoadAllService() error {
	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/services/")

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	resp, err := self.client.Get(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for k := range resp.Kvs {
		err = self.LoadService(global.BytesToStr(resp.Kvs[k].Key), resp.Kvs[k].Value)
		if err != nil {
			return err
		}
	}
	return nil
}

// 将本地服务数据保存到存储器中，如果不存在则创建
func (self *Etcd) SaveService(serviceID, loadBalance string) (err error) {
	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/services/")
	key.WriteString(global.EncodeKey(serviceID))

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err = self.client.Put(ctx, key.String(), loadBalance); err != nil {
		return err
	}
	return nil
}

// 将本地中所有服务数据保存到存储器
func (self *Etcd) SaveAllService() error {
	var (
		err      error
		key      strings.Builder
		services = make(map[string]string, global.SyncMapLen(&global.Services))
	)

	// 将配置保存到临时变量中
	global.Services.Range(func(k, v interface{}) bool {
		service, ok := v.(global.ServiceType)
		if !ok {
			err = errors.New("服务" + k.(string) + "的配置异常")
			log.Err(err).Caller().Msg("类型断言失败")
			return false
		}
		services[service.ID] = service.LoadBalance
		return true
	})
	if err != nil {
		return err
	}

	// 清空存储器中的配置
	key.WriteString(self.KeyPrefix)
	key.WriteString("/services/")
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	_, err = self.client.Delete(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		log.Err(err).Caller().Msg("清空存储器中的数据失败")
		return err
	}

	// 将内存中的数据写入到存储器中
	for serviceID, loadBalance := range services {
		if err = self.SaveService(serviceID, loadBalance); err != nil {
			return err
		}
	}

	return nil
}

// 删除本地服务数据
func (self *Etcd) DeleteLocalService(key string) error {
	serviceID, err := global.DecodeKey(path.Base(key))
	if err != nil {
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
		log.Err(err).Caller().Msg("删除存储器中的服务数据失败")
		return err
	}
	return nil
}
