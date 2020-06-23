package etcd

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/rs/zerolog/log"

	"local/global"
)

// 加载所有数据
// todo 加载所有数据要加分布式锁，防止加载的不是最新的数据
func (self *Etcd) LoadAll() error {
	// 清空本地数据
	global.Services.Range(func(key, _ interface{}) bool {
		global.Services.Delete(key)
		return true
	})
	global.Services = sync.Map{}
	atomic.StoreUint32(&global.TotalServices, 0)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	// 从远程加载所有服务列表
	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/services/")
	resp, err := self.client.Get(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		log.Err(err).Caller().Send()
		return err
	}
	for k := range resp.Kvs {
		err = self.LoadService(resp.Kvs[k].Value)
		if err != nil {
			log.Err(err).Caller().Send()
			return err
		}
	}

	// 从远程加载所有节点列表
	key.Reset()
	key.WriteString(self.KeyPrefix)
	key.WriteString("/nodes/")
	resp, err = self.client.Get(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		log.Err(err).Caller().Send()
		return err
	}
	for k := range resp.Kvs {
		err = self.LoadNode(global.BytesToStr(resp.Kvs[k].Key), resp.Kvs[k].Value)
		if err != nil {
			log.Err(err).Caller().Send()
			return err
		}
	}
	return nil
}

// 保存所有数据到远程
// todo 保存所有数据时要加分布式锁，防止出现原子性问题
func (self *Etcd) SaveAll() (err error) {
	global.Services.Range(func(key, value interface{}) bool {
		ci, ok := value.(global.Cluster)
		if !ok {
			log.Error().Caller().Send()
			return false
		}
		if err = self.SaveService(ci.Config()); err != nil {
			log.Err(err).Caller().Send()
			return false
		}
		nodes := ci.Nodes()
		for k := range nodes {
			if err = self.SaveNode(ci.Config().ServiceID, nodes[k].IP, nodes[k].Port, nodes[k].Weight, nodes[k].Expires); err != nil {
				log.Err(err).Caller().Send()
				return false
			}
		}
		return true
	})
	return
}
