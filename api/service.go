package api

import (
	"net/http"

	"local/engine"
	"local/global"

	"github.com/dxvgef/filter/v2"
	"github.com/dxvgef/tsing"
)

type Service struct{}

func (self *Service) Add(ctx *tsing.Context) error {
	var (
		err    error
		resp   = make(map[string]string)
		config global.ServiceConfig
	)
	if err = filter.Batch(
		filter.String(ctx.Post("id"), "id").Require().Set(&config.ServiceID),
		filter.String(ctx.Post("load_balance"), "load_balance").Require().Set(&config.LoadBalance),
	); err != nil {
		// 来自客户端的数据，无需记录日志
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	if _, exists := global.Services.Load(config.ServiceID); exists {
		resp["error"] = "服务ID已存在"
		return JSON(ctx, 400, &resp)
	}
	if config.LoadBalance == "" {
		resp["error"] = "load_balance参数不能为空"
		return JSON(ctx, 400, &resp)
	}

	if err = global.Storage.SaveService(config); err != nil {
		return ctx.Caller(err)
	}

	return Status(ctx, 204)
}
func (self *Service) Put(ctx *tsing.Context) error {
	var (
		err    error
		resp   = make(map[string]string)
		config global.ServiceConfig
	)
	if err = filter.Batch(
		filter.String(ctx.PathParams.Value("serviceID"), "serviceID").Require().Base64RawURLDecode().Set(&config.ServiceID),
		filter.String(ctx.Post("load_balance"), "load_balance").Require().Set(&config.LoadBalance),
	); err != nil {
		// 来自客户端的数据，无需记录日志
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	if config.LoadBalance == "" {
		resp["error"] = "load_balance参数不能为空"
		return JSON(ctx, 400, &resp)
	}

	if err = global.Storage.SaveService(config); err != nil {
		return ctx.Caller(err)
	}

	return Status(ctx, 204)
}

func (self *Service) Delete(ctx *tsing.Context) error {
	var (
		err       error
		serviceID string
	)
	if serviceID, err = global.DecodeKey(ctx.PathParams.Value("serviceID")); err != nil {
		// 来自客户端的数据，无需记录日志
		return Status(ctx, 404)
	}
	if _, exist := global.Services.Load(serviceID); !exist {
		return Status(ctx, 404)
	}
	err = global.Storage.DeleteStorageService(ctx.PathParams.Value("serviceID"))
	if err != nil {
		return ctx.Caller(err)
	}
	return Status(ctx, 204)
}

// 选取节点
func (self *Service) Select(ctx *tsing.Context) error {
	var (
		err       error
		resp      = make(map[string]string)
		serviceID string
	)
	if serviceID, err = filter.String(ctx.PathParams.Value("serviceID"), "serviceID").Require().Base64RawURLDecode().String(); err != nil {
		// 来自客户端的数据，无需记录日志
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	ci := engine.FindCluster(serviceID)
	if ci == nil {
		resp["error"] = "服务不存在"
		return JSON(ctx, 400, &resp)
	}
	node := ci.Select()
	if resp["ip"] == "" {
		return Status(ctx, http.StatusNotImplemented)
	}
	return JSON(ctx, 200, &node)
}
