package etcd

import (
	"context"
	"errors"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/dxvgef/tsing-center/engine"
	"github.com/dxvgef/tsing-center/global"

	"github.com/coreos/etcd/clientv3"
	"github.com/rs/zerolog/log"
)

// 从存储器加载节点到本地，如果不存在则创建
func (self *Etcd) LoadNode(key string, data []byte) error {
	var keyPrefix strings.Builder
	keyPrefix.WriteString(self.KeyPrefix)
	keyPrefix.WriteString("/nodes/")
	serviceID, ip, port, err := self.ParseNode(key, keyPrefix.String())
	if err != nil {
		return err
	}
	var weight int
	weight, err = strconv.Atoi(global.BytesToStr(data))
	if err != nil {
		return err
	}
	return engine.SetNode(serviceID, ip, port, weight)
}

// 从存储器加载所有节点数据到本地
func (self *Etcd) LoadAllNode() error {
	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/nodes/")

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

// 将本地节点数据保存到存储器中，如果不存在则创建
func (self *Etcd) SaveNode(serviceID, ip string, port uint16, weight int) (err error) {
	var key strings.Builder
	key.WriteString(ip)
	key.WriteString(":")
	key.WriteString(strconv.FormatUint(uint64(port), 10))
	node := global.EncodeKey(key.String())

	key.Reset()
	key.WriteString(self.KeyPrefix)
	key.WriteString("/nodes/")
	key.WriteString(global.EncodeKey(serviceID))
	key.WriteString("/")
	key.WriteString(node)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err = self.client.Put(ctx, key.String(), strconv.Itoa(weight)); err != nil {
		return err
	}
	return nil
}

// 将本地中所有节点数据保存到存储器
func (self *Etcd) SaveAllNode() error {
	var (
		err error
		key strings.Builder
		lbs = make(map[string]global.LoadBalance, global.SyncMapLen(&global.Nodes))
	)

	// 将配置保存到临时变量中
	global.Nodes.Range(func(k, v interface{}) bool {
		node, ok := v.(global.LoadBalance)
		if !ok {
			err = errors.New("断言失败")
			log.Err(err).Caller().Msg("类型断言失败")
			return false
		}
		lbs[k.(string)] = node
		return true
	})
	if err != nil {
		return err
	}

	// 清空存储器中的配置
	key.WriteString(self.KeyPrefix)
	key.WriteString("/lbs/")
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	_, err = self.client.Delete(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		log.Err(err).Caller().Msg("清空存储器中的数据失败")
		return err
	}

	// 将内存中的数据写入到存储器中
	for serviceID, lb := range lbs {
		nodes := lb.Nodes()
		for k := range nodes {
			if err = self.SaveNode(serviceID, nodes[k].IP, nodes[k].Port, nodes[k].Weight); err != nil {
				return err
			}
		}
	}

	return nil
}

// 删除本地的节点
func (self *Etcd) DeleteLocalNode(key string) error {
	var keyPrefix strings.Builder
	keyPrefix.WriteString(self.KeyPrefix)
	keyPrefix.WriteString("/nodes/")
	serviceID, ip, port, err := self.ParseNode(key, keyPrefix.String())
	if err != nil {
		return err
	}
	return engine.DelNode(serviceID, ip, port)
}

// 删除存储器的节点
func (self *Etcd) DeleteStorageNode(serviceID, ip string, port uint16) error {
	if serviceID == "" {
		return errors.New("serviceID不能为空")
	}
	if ip == "" {
		return errors.New("ip不能为空")
	}
	if port == 0 {
		return errors.New("port不能为空")
	}
	var key strings.Builder
	key.WriteString(ip)
	key.WriteString(":")
	key.WriteString(strconv.FormatUint(uint64(port), 10))
	node := global.EncodeKey(key.String())

	key.Reset()
	key.WriteString(self.KeyPrefix)
	key.WriteString("/nodes/")
	key.WriteString(global.EncodeKey(serviceID))
	key.WriteString("/")
	key.WriteString(node)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	_, err := self.client.Delete(ctx, key.String())
	if err != nil {
		log.Err(err).Caller().Msg("删除存储器中的服务数据失败")
		return err
	}
	return nil
}

// 从key字符串中解析节点信息
// key=prefix/nodes/base64(serviceID)/base64(ip:port)
func (self *Etcd) ParseNode(key, prefix string) (serviceID, ip string, port uint16, err error) {
	if prefix != "" {
		key = strings.TrimPrefix(key, prefix)
	}

	var nodePart string
	nodePart, err = global.DecodeKey(path.Base(key))
	if err != nil {
		return
	}
	pos := strings.Index(nodePart, ":")
	if pos == -1 {
		err = errors.New("解析节点信息失败")
		return
	}
	ip = nodePart[0:pos]
	p, er := strconv.Atoi(nodePart[pos:])
	if er != nil {
		err = er
		return
	}
	port = uint16(p)

	pos = -1
	pos = strings.Index(key, "/")
	if pos == -1 {
		err = errors.New("解析服务ID信息失败")
		return
	}

	serviceID, err = global.DecodeKey(key[0:pos])
	return
}
