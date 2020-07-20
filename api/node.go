package api

import (
	"math"
	"strconv"
	"strings"

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
			expires   int64
		}
	)
	if err = filter.Batch(
		filter.String(ctx.Post("service_id"), "service_id").Require().Set(&req.serviceID),
		filter.String(ctx.Post("ip"), "ip").Require().IsIP().Set(&req.ip),
		filter.String(ctx.Post("port"), "port").Require().IsDigit().MinInteger(1).MaxInteger(math.MaxUint16).Set(&req.port),
		filter.String(ctx.Post("weight"), "weight").Require().MinInteger(0).MaxInteger(math.MaxUint16).Set(&req.weight),
		filter.String(ctx.Post("expires"), "expires").MinInteger(0).IsDigit().Set(&req.expires),
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

	if err = global.Storage.SaveNode(req.serviceID, req.ip, req.port, req.weight, req.expires); err != nil {
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
			expires   int64
		}
		port uint64
	)
	if err = filter.Batch(
		filter.String(ctx.PathParams.Value("serviceID"), "serviceID").Require().Base64RawURLDecode().Set(&req.serviceID),
		filter.String(ctx.PathParams.Value("node"), "node").Require().Base64RawURLDecode().Set(&req.node),
		filter.String(ctx.Post("weight"), "weight").Require().MinInteger(0).MaxInteger(math.MaxUint16).Set(&req.weight),
		filter.String(ctx.Post("expires"), "expires").MinInteger(0).IsDigit().Set(&req.expires),
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
	if err = global.Storage.SaveNode(req.serviceID, req.ip, req.port, req.weight, req.expires); err != nil {
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
			ip        string
			port      uint16
		}
		port uint64
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
	req.ip = req.node[0:pos]
	port, err = strconv.ParseUint(req.node[pos:], 10, 16)
	if err != nil {
		// 来自客户端的数据，无需记录日志
		return Status(ctx, 404)
	}
	req.port = uint16(port)

	err = global.Storage.DeleteStorageNode(req.serviceID, req.ip, req.port)
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
		}
		port64 uint64
	)

	// 验证请求参数
	if err = filter.Batch(
		filter.String(ctx.PathParams.Value("serviceID"), "serviceID").Require().Base64RawURLDecode().Set(&req.serviceID),
		filter.String(ctx.PathParams.Value("node"), "node").Require().Base64RawURLDecode().Set(&req.node),
		filter.String(ctx.PathParams.Value("attrs"), "attrs").Require().EnumSliceString(",", []string{"expires", "weight"}).SetSlice(&req.attrs, ","),
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
		if req.attrs[k] == "expires" {
			if _, exist := ctx.PostParam("expires"); !exist {
				resp["error"] = "expires值不能为空"
				return JSON(ctx, 400, &resp)
			}
			node.Expires, err = filter.String(ctx.Post("expires"), "expires").Require().IsDigit().MinInteger(0).Int64()
			if err != nil {
				resp["error"] = err.Error()
				return JSON(ctx, 400, &resp)
			}
		}
	}

	// 更新本地内存数据
	ci.Set(node.IP, node.Port, node.Weight, node.Expires)

	// 更新存储引擎中的数据
	if err = global.Storage.SaveNode(req.serviceID, node.IP, node.Port, node.Weight, node.Expires); err != nil {
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
			expires   int64
		}
		port64 uint64
	)

	// 验证请求参数
	if err = filter.Batch(
		filter.String(ctx.PathParams.Value("serviceID"), "serviceID").Require().Base64RawURLDecode().Set(&req.serviceID),
		filter.String(ctx.PathParams.Value("node"), "node").Require().Base64RawURLDecode().Set(&req.node),
		filter.String(ctx.Post("expires"), "expires").Require().IsDigit().MinInteger(0).Set(&req.expires),
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

	// 更新本地内存数据
	ci.Set(req.ip, req.port, node.Weight, req.expires)

	// 更新存储引擎中的数据
	if err = global.Storage.SaveNode(req.serviceID, node.IP, node.Port, node.Weight, req.expires); err != nil {
		return ctx.Caller(err)
	}

	return Status(ctx, 204)
}
