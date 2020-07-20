package api

import (
	"math"
	"strconv"
	"strings"
	"time"

	"local/engine"
	"local/global"

	"github.com/dxvgef/filter/v2"
	"github.com/dxvgef/tsing"
)

type Node struct{}

func (self *Node) Add(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
		req  struct {
			serviceID string
			ip        string
			port      uint16
			weight    int
			ttl       uint
			expires   int64
			meta      string
		}
	)
	if err = filter.Batch(
		filter.String(ctx.Post("service_id"), "service_id").Require().Set(&req.serviceID),
		filter.String(ctx.Post("ip"), "ip").Require().IsIP().Set(&req.ip),
		filter.String(ctx.Post("port"), "port").Require().IsDigit().MinInteger(1).MaxInteger(math.MaxUint16).Set(&req.port),
		filter.String(ctx.Post("weight"), "weight").Require().MinInteger(0).MaxInteger(math.MaxUint16).Set(&req.weight),
		filter.String(ctx.Post("ttl"), "ttl").MinInteger(0).IsDigit().Set(&req.ttl),
		filter.String(ctx.Post("meta"), "meta").IsJSON().Set(&req.meta),
	); err != nil {
		// 来自客户端的数据，无需记录日志
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}

	// 自动获取节点IP
	if req.ip == "" {
		req.ip = ctx.Request.RemoteAddr[:strings.Index(ctx.Request.RemoteAddr, ":")]
	}

	ci := engine.FindCluster(req.serviceID)
	if ci == nil {
		resp["error"] = "服务不存在"
		return JSON(ctx, 400, &resp)
	}
	node := ci.Find(req.ip, req.port)
	if node.IP != "" {
		resp["error"] = "节点已存在"
		return JSON(ctx, 400, &resp)
	}

	if req.ttl > 0 {
		req.expires = time.Now().Add(time.Duration(req.ttl) * time.Second).Unix()
	}

	if err = global.Storage.SaveNode(req.serviceID, global.Node{
		IP:      req.ip,
		Port:    req.port,
		Weight:  req.weight,
		TTL:     req.ttl,
		Expires: req.expires,
		Mete:    req.meta,
	}); err != nil {
		return ctx.Caller(err)
	}

	return Status(ctx, 204)
}

func (self *Node) Put(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
		req  struct {
			serviceID string
			node      string
			ip        string
			port      uint16
			weight    int
			ttl       uint
			expires   int64
			meta      string
		}
		port uint64
	)
	if err = filter.Batch(
		filter.String(ctx.PathParams.Value("serviceID"), "serviceID").Require().Base64RawURLDecode().Set(&req.serviceID),
		filter.String(ctx.PathParams.Value("node"), "node").Require().Base64RawURLDecode().Set(&req.node),
		filter.String(ctx.Post("weight"), "weight").Require().MinInteger(0).MaxInteger(math.MaxUint16).Set(&req.weight),
		filter.String(ctx.Post("ttl"), "ttl").MinInteger(0).IsDigit().Set(&req.ttl),
		filter.String(ctx.Post("meta"), "meta").IsJSON().Set(&req.meta),
	); err != nil {
		// 来自客户端的数据，无需记录日志
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	pos := strings.Index(req.node, ":")
	if pos == -1 {
		// 来自客户端的数据，无需记录日志
		return Status(ctx, 404)
	}
	req.ip = req.node[0:pos]
	port, err = strconv.ParseUint(req.node[pos+1:], 10, 16)
	if err != nil {
		// 来自客户端的数据，无需记录日志
		return Status(ctx, 404)
	}
	req.port = uint16(port)

	ci := engine.FindCluster(req.serviceID)
	if ci == nil {
		resp["error"] = "服务不存在"
		return JSON(ctx, 400, &resp)
	}

	if req.ttl > 0 {
		req.expires = time.Now().Add(time.Duration(req.ttl) * time.Second).Unix()
	}

	if err = global.Storage.SaveNode(req.serviceID, global.Node{
		IP:      req.ip,
		Port:    req.port,
		Weight:  req.weight,
		TTL:     req.ttl,
		Expires: req.expires,
		Mete:    req.meta,
	}); err != nil {
		return ctx.Caller(err)
	}

	return Status(ctx, 204)
}

func (self *Node) Delete(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
		req  struct {
			serviceID string
			node      string
		}
		ip     string
		port   uint16
		port64 uint64
	)
	if err = filter.Batch(
		filter.String(ctx.PathParams.Value("serviceID"), "serviceID").Require().Base64RawURLDecode().Set(&req.serviceID),
		filter.String(ctx.PathParams.Value("node"), "node").Require().Base64RawURLDecode().Set(&req.node),
	); err != nil {
		// 来自客户端的数据，无需记录日志
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	pos := strings.Index(req.node, ":")
	if pos == -1 {
		// 来自客户端的数据，无需记录日志
		return Status(ctx, 404)
	}
	ip = req.node[0:pos]
	port64, err = strconv.ParseUint(req.node[pos+1:], 10, 16)
	if err != nil {
		// 来自客户端的数据，无需记录日志
		return Status(ctx, 404)
	}
	port = uint16(port64)

	err = global.Storage.DeleteStorageNode(req.serviceID, ip, port)
	if err != nil {
		return ctx.Caller(err)
	}
	return Status(ctx, 204)
}

// 更新节点属性
func (self *Node) Patch(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
		req  struct {
			serviceID string
			node      string
			ip        string
			port      uint16
			attrs     []string
			expires   int64
			meta      string
		}
		port64 uint64
	)

	// 验证请求参数
	if err = filter.Batch(
		filter.String(ctx.PathParams.Value("serviceID"), "serviceID").Require().Base64RawURLDecode().Set(&req.serviceID),
		filter.String(ctx.PathParams.Value("node"), "node").Require().Base64RawURLDecode().Set(&req.node),
		filter.String(ctx.PathParams.Value("attrs"), "attrs").Require().EnumSliceString(",", []string{"ttl", "weight", "meta"}).SetSlice(&req.attrs, ","),
	); err != nil {
		// 来自客户端的数据，无需记录日志
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	pos := strings.Index(req.node, ":")
	if pos == -1 {
		// 来自客户端的数据，无需记录日志
		return Status(ctx, 404)
	}
	req.ip = req.node[0:pos]
	port64, err = strconv.ParseUint(req.node[pos+1:], 10, 16)
	if err != nil {
		// 来自客户端的数据，无需记录日志
		return Status(ctx, 404)
	}
	req.port = uint16(port64)

	// 获取集群
	ci := engine.FindCluster(req.serviceID)
	if ci == nil {
		// 来自客户端的数据，无需记录日志
		return Status(ctx, 404)
	}

	// 获取节点
	node := ci.Find(req.ip, req.port)
	if node.IP == "" {
		return Status(ctx, 404)
	}

	for k := range req.attrs {
		if req.attrs[k] == "weight" {
			if _, exist := ctx.PostParam("weight"); !exist {
				resp["error"] = "weight值无效"
				return JSON(ctx, 400, &resp)
			}
			node.Weight, err = filter.String(ctx.Post("weight"), "weight").Require().IsDigit().MinInteger(0).Int()
			if err != nil {
				resp["error"] = err.Error()
				return JSON(ctx, 400, &resp)
			}
		}
		if req.attrs[k] == "meta" {
			node.Mete, err = filter.String(ctx.Post("weight"), "weight").IsJSON().String()
			if err != nil {
				resp["error"] = err.Error()
				return JSON(ctx, 400, &resp)
			}
		}
		if req.attrs[k] == "ttl" {
			if _, exist := ctx.PostParam("ttl"); !exist {
				resp["error"] = "ttl值无效"
				return JSON(ctx, 400, &resp)
			}
			node.TTL, err = filter.String(ctx.Post("ttl"), "ttl").Require().IsDigit().MinInteger(0).Uint()
			if err != nil {
				resp["error"] = err.Error()
				return JSON(ctx, 400, &resp)
			}
			if node.TTL != 0 {
				node.Expires = time.Now().Add(time.Duration(node.TTL) * time.Second).Unix()
			} else {
				node.Expires = 0
			}
		}
	}

	// 更新存储引擎中的数据
	if err = global.Storage.SaveNode(req.serviceID, global.Node{
		IP:      node.IP,
		Port:    node.Port,
		Weight:  node.Weight,
		TTL:     node.TTL,
		Mete:    node.Mete,
		Expires: node.Expires,
	}); err != nil {
		return ctx.Caller(err)
	}

	return Status(ctx, 204)
}

// 更新节点生命周期的截止时间
func (self *Node) Touch(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
		req  struct {
			serviceID string
			node      string
			ip        string
			port      uint16
			ttl       uint
		}
		port64 uint64
	)

	// 验证请求参数
	if err = filter.Batch(
		filter.String(ctx.PathParams.Value("serviceID"), "serviceID").Require().Base64RawURLDecode().Set(&req.serviceID),
		filter.String(ctx.PathParams.Value("node"), "node").Require().Base64RawURLDecode().Set(&req.node),
	); err != nil {
		// 来自客户端的数据，无需记录日志
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	pos := strings.Index(req.node, ":")
	if pos == -1 {
		// 来自客户端的数据，无需记录日志
		return Status(ctx, 404)
	}
	req.ip = req.node[0:pos]
	port64, err = strconv.ParseUint(req.node[pos+1:], 10, 16)
	if err != nil {
		// 来自客户端的数据，无需记录日志
		return Status(ctx, 404)
	}
	req.port = uint16(port64)

	// 获取集群
	ci := engine.FindCluster(req.serviceID)
	if ci == nil {
		// 来自客户端的数据，无需记录日志
		return Status(ctx, 404)
	}

	// 获取节点
	node := ci.Find(req.ip, req.port)
	if node.IP == "" {
		return Status(ctx, 404)
	}

	if node.TTL == 0 {
		return Status(ctx, 204)
	}

	// 更新存储引擎中的数据
	expires := time.Now().Add(time.Duration(req.ttl) * time.Second).Unix()
	if err = global.Storage.SaveNode(req.serviceID, global.Node{
		IP:      node.IP,
		Port:    node.Port,
		Weight:  node.Weight,
		TTL:     node.TTL,
		Expires: expires,
		Mete:    node.Mete,
	}); err != nil {
		return ctx.Caller(err)
	}

	return Status(ctx, 204)
}
