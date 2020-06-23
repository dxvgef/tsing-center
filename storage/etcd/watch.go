package etcd

import (
	"context"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/rs/zerolog/log"

	"local/global"
)

// 监听变更
func (self *Etcd) Watch() {
	ch := self.client.Watch(context.Background(), self.KeyPrefix+"/", clientv3.WithPrefix())
	for resp := range ch {
		for _, event := range resp.Events {
			switch event.Type {
			// 更新事件
			case clientv3.EventTypePut:
				if err := self.watchLoadData(event.Kv.Key, event.Kv.Value); err != nil {
					log.Err(err).Caller().Send()
				}
			// 删除事件
			case clientv3.EventTypeDelete:
				if err := self.watchDeleteData(event.Kv.Key); err != nil {
					log.Err(err).Caller().Send()
				}
			}
		}
	}
}

// 监听存储器数据更新，同步本地数据
func (self *Etcd) watchLoadData(key, value []byte) error {
	keyStr := global.BytesToStr(key)
	// 加载服务
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/services/") {
		return self.LoadService(value)
	}
	// 加载节点
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/nodes/") {
		return self.LoadNode(keyStr, value)
	}
	return nil
}

// 监听存储器数据删除，同步本地数据
func (self *Etcd) watchDeleteData(key []byte) error {
	keyStr := global.BytesToStr(key)
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/services/") {
		return self.DeleteLocalService(keyStr)
	}
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/nodes/") {
		return self.DeleteLocalNode(keyStr)
	}
	return nil
}
