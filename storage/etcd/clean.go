package etcd

import (
	"github.com/rs/zerolog/log"

	"local/global"
)

// 批理清理无效的节点
func (self *Etcd) Clean(serviceID string, nodes []global.Node) (err error) {
	for k := range nodes {
		if err = self.DeleteStorageNode(serviceID, nodes[k].IP, nodes[k].Port); err != nil {
			log.Err(err).Str("ip", nodes[k].IP).Uint16("port", nodes[k].Port).Caller().Send()
			return
		}
	}
	return
}
