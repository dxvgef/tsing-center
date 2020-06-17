package api

import (
	"net/http"

	"github.com/dxvgef/tsing-center/engine"
	"github.com/dxvgef/tsing-center/global"

	"github.com/dxvgef/filter"
	"github.com/dxvgef/tsing"
)

type Service struct{}

func (self *Service) Add(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
		req  struct {
			serviceID   string
			loadBalance string
		}
	)
	if err = filter.MSet(
		filter.El(&req.serviceID, filter.FromString(ctx.Post("id"), "id").Required()),
		filter.El(&req.loadBalance, filter.FromString(ctx.Post("load_balance"), "load_balance")),
	); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	if _, exists := global.Services.Load(req.serviceID); exists {
		resp["error"] = "服务ID已存在"
		return JSON(ctx, 400, &resp)
	}
	if req.loadBalance == "" {
		resp["error"] = "load_balance参数不能为空"
		return JSON(ctx, 400, &resp)
	}

	if err = global.Storage.SaveService(req.serviceID, req.loadBalance); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}

	return Status(ctx, 204)
}
func (self *Service) Put(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
		req  struct {
			serviceID   string
			loadBalance string
		}
	)
	if err = filter.MSet(
		filter.El(&req.serviceID, filter.FromString(ctx.PathParams.Value("serviceID"), "serviceID").Required().Base64RawURLDecode()),
		filter.El(&req.loadBalance, filter.FromString(ctx.Post("load_balance"), "load_balance")),
	); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	if req.loadBalance == "" {
		resp["error"] = "load_balance参数不能为空"
		return JSON(ctx, 400, &resp)
	}

	if err = global.Storage.SaveService(req.serviceID, req.loadBalance); err != nil {
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
	lb, exist := engine.MatchNode(serviceID)
	if !exist {
		resp["error"] = "服务不存在"
		return JSON(ctx, 400, &resp)
	}
	ip, port := lb.Next()
	if ip == "" {
		return Status(ctx, http.StatusNotImplemented)
	}
	resp["ip"] = ip
	resp["port"] = port
	return JSON(ctx, 200, &resp)
}
