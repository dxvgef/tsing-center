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

	"github.com/rs/zerolog/log"
)

// 节点值结构
type NodeValue struct {
	Expires int64 `json:"expires,omitempty"` // 生命周期截止时间(unix时间戳)
	Weight  int   `json:"weight,omitempty"`  // 权重值
}

// 从存储器加载节点到本地，如果不存在则创建
func (self *Etcd) LoadNode(key string, data []byte) error {
	var keyPrefix strings.Builder
	keyPrefix.WriteString(self.KeyPrefix)
	keyPrefix.WriteString("/nodes/")

	// 从key中解析serviceID, ip, port
	serviceID, ip, port, err := self.ParseNode(key, keyPrefix.String())
	if err != nil {
		return err
	}

	// 从value中解析expires, weight
	var value NodeValue
	err = value.UnmarshalJSON(data)
	if err != nil {
		log.Err(err).Caller().Msg("解码节点属性失败")
		return err
	}

	// 写入节点到本地
	return engine.SetNode(serviceID, ip, port, value.Weight, value.Expires)
}

// 将本地节点数据保存到存储器中，如果不存在则创建
func (self *Etcd) SaveNode(serviceID, ip string, port uint16, weight int, expires int64) (err error) {
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

	var (
		value      NodeValue
		valueBytes []byte
	)
	value.Weight = weight
	value.Expires = expires
	valueBytes, err = value.MarshalJSON()
	if err != nil {
		log.Err(err).Caller().Msg("编码节点属性失败")
		return
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err = self.client.Put(ctx, key.String(), global.BytesToStr(valueBytes)); err != nil {
		return err
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
		log.Error().Caller().Msg("解析节点信息失败")
		return
	}
	ip = nodePart[0:pos]
	p, er := strconv.Atoi(nodePart[pos+1:])
	if er != nil {
		err = er
		log.Err(err).Caller().Msg("解析节点信息失败")
		return
	}
	port = uint16(p)

	pos = -1
	pos = strings.Index(key, "/")
	if pos == -1 {
		err = errors.New("解析服务ID信息失败")
		log.Error().Caller().Msg("解析服务ID信息失败")
		return
	}

	serviceID, err = global.DecodeKey(key[0:pos])
	return
}
