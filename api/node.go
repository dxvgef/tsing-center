package api

import (
	"math"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

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
		log.Err(err).Caller().Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
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
		log.Err(err).Caller().Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
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
		log.Err(err).Caller().Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

// 更新到期时间
func (self *Node) UpdateExpires(ctx *tsing.Context) error {
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

	ci := engine.FindCluster(req.serviceID)
	if ci == nil {
		// 来自客户端的数据，无需记录日志
		return Status(ctx, 404)
	}

	ci.Set(req.ip, req.port, -1, req.expires)
	return Status(ctx, 204)
}
