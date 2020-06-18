package api

import (
	"net/http"
	"time"

	"github.com/dxvgef/tsing-center/engine"
	"github.com/dxvgef/tsing-center/global"

	"github.com/dxvgef/filter"
	"github.com/dxvgef/tsing"
)

type Service struct{}

func (self *Service) Add(ctx *tsing.Context) error {
	var (
		err    error
		resp   = make(map[string]string)
		config global.ServiceConfig
	)
	if err = filter.MSet(
		filter.El(&config.ServiceID, filter.FromString(ctx.Post("id"), "id").Required()),
		filter.El(&config.LoadBalance, filter.FromString(ctx.Post("load_balance"), "load_balance")),
	); err != nil {
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
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}

	return Status(ctx, 204)
}
func (self *Service) Put(ctx *tsing.Context) error {
	var (
		err    error
		resp   = make(map[string]string)
		config global.ServiceConfig
	)
	if err = filter.MSet(
		filter.El(&config.ServiceID, filter.FromString(ctx.PathParams.Value("serviceID"), "serviceID").Required().Base64RawURLDecode()),
		filter.El(&config.LoadBalance, filter.FromString(ctx.Post("load_balance"), "load_balance")),
	); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	if config.LoadBalance == "" {
		resp["error"] = "load_balance参数不能为空"
		return JSON(ctx, 400, &resp)
	}

	if err = global.Storage.SaveService(config); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}

	return Status(ctx, 204)
}

func (self *Service) Delete(ctx *tsing.Context) error {
	var (
		err       error
		resp      = make(map[string]string)
		serviceID string
	)
	if serviceID, err = global.DecodeKey(ctx.PathParams.Value("serviceID")); err != nil {
		return Status(ctx, 404)
	}
	if _, exist := global.Services.Load(serviceID); !exist {
		return Status(ctx, 404)
	}
	err = global.Storage.DeleteStorageService(ctx.PathParams.Value("serviceID"))
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Service) Next(ctx *tsing.Context) error {
	var (
		err       error
		resp      = make(map[string]interface{})
		serviceID string
	)
	if err = filter.MSet(
		filter.El(&serviceID, filter.FromString(ctx.PathParams.Value("serviceID"), "serviceID").Required().Base64RawURLDecode()),
	); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	ci := engine.FindCluster(serviceID)
	if ci == nil {
		resp["error"] = "服务不存在"
		return JSON(ctx, 400, &resp)
	}
	// 最多三次重新选取节点的机会
	for i := 0; i < 3; i++ {
		ip, port, expires := ci.Next()
		// 如果存在有效期，并且已过期
		if expires > 0 && expires < time.Now().Unix() {
			// 删除该节点
			ci.Remove(ip, port)
			continue
		}
		resp["ip"] = ip
		resp["port"] = port
		break
	}
	if resp["ip"] == "" {
		return Status(ctx, http.StatusNotImplemented)
	}
	return JSON(ctx, 200, &resp)
}
