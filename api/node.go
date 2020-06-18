package api

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-center/engine"
	"github.com/dxvgef/tsing-center/global"

	"github.com/dxvgef/filter"
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
	if err = filter.MSet(
		filter.El(&req.serviceID, filter.FromString(ctx.Post("service_id"), "service_id").Required()),
		filter.El(&req.ip, filter.FromString(ctx.Post("ip"), "ip").Required().IsIP()),
		filter.El(&req.port, filter.FromString(ctx.Post("port"), "port").Required().IsDigit().MinInteger(1).MaxInteger(math.MaxUint16)),
		filter.El(&req.weight, filter.FromString(ctx.Post("weight"), "weight").Required().MinInteger(0).MaxInteger(math.MaxUint16)),
		filter.El(&req.expires, filter.FromString(ctx.Post("expires"), "expires").MinInteger(0).IsDigit()),
	); err != nil {
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
	if err = filter.MSet(
		filter.El(&req.serviceID, filter.FromString(ctx.PathParams.Value("serviceID"), "serviceID").Required().Base64RawURLDecode()),
		filter.El(&req.node, filter.FromString(ctx.PathParams.Value("node"), "node").Required().Base64RawURLDecode()),
		filter.El(&req.weight, filter.FromString(ctx.Post("weight"), "weight").Required().MinInteger(0).MaxInteger(math.MaxUint16)),
		filter.El(&req.expires, filter.FromString(ctx.Post("expires"), "expires").MinInteger(0).IsDigit()),
	); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	pos := strings.Index(req.node, ":")
	if pos == -1 {
		log.Debug().Int("pos", pos).Msg("解析node失败")
		return Status(ctx, 404)
	}
	req.ip = req.node[0:pos]
	port, err = strconv.ParseUint(req.node[pos+1:], 10, 16)
	if err != nil {
		return Status(ctx, 404)
	}
	req.port = uint16(port)

	ci := engine.FindCluster(req.serviceID)
	if ci == nil {
		resp["error"] = "服务不存在"
		return JSON(ctx, 400, &resp)
	}
	if err = global.Storage.SaveNode(req.serviceID, req.ip, req.port, req.weight, req.expires); err != nil {
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
	if err = filter.MSet(
		filter.El(&req.serviceID, filter.FromString(ctx.PathParams.Value("serviceID"), "serviceID").Required().Base64RawURLDecode()),
		filter.El(&req.node, filter.FromString(ctx.PathParams.Value("node"), "node").Required().Base64RawURLDecode()),
	); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	pos := strings.Index(req.node, ":")
	if pos == -1 {
		return Status(ctx, 404)
	}
	req.ip = req.node[0:pos]
	port, err = strconv.ParseUint(req.node[pos:], 10, 16)
	if err != nil {
		return Status(ctx, 404)
	}
	req.port = uint16(port)

	err = global.Storage.DeleteStorageNode(req.serviceID, req.ip, req.port)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

// 促活
func (self *Node) Active(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
		req  struct {
			serviceID string
			node      string
			ip        string
			port      uint16
		}
		port64 uint64
	)
	if err = filter.MSet(
		filter.El(&req.serviceID, filter.FromString(ctx.PathParams.Value("serviceID"), "serviceID").Required().Base64RawURLDecode()),
		filter.El(&req.node, filter.FromString(ctx.PathParams.Value("node"), "node").Required().Base64RawURLDecode()),
	); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	pos := strings.Index(req.node, ":")
	if pos == -1 {
		log.Debug().Int("pos", pos).Msg("解析node失败")
		return Status(ctx, 404)
	}
	req.ip = req.node[0:pos]
	port64, err = strconv.ParseUint(req.node[pos+1:], 10, 16)
	if err != nil {
		return Status(ctx, 404)
	}
	req.port = uint16(port64)

	ci := engine.FindCluster(req.serviceID)
	if ci == nil {
		return Status(ctx, 404)
	}

	ci.Set(req.ip, req.port, -1, time.Now().Add(10*time.Second).Unix())
	return Status(ctx, 204)
}
