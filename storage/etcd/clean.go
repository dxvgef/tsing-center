package etcd

import (
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-center/global"
)

// 批理清理无效的节点
func (self *Etcd) Clean(nodes map[string]global.Node) (err error) {
	for k, v := range nodes {
		if err = global.Storage.DeleteStorageNode(k, v.IP, v.Port); err != nil {
			log.Err(err).Caller().Msg("清理无效的节点失败")
			return
		}
	}
	return
}
