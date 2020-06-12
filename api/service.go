package api

import (
	"github.com/dxvgef/filter"
	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-center/global"
)

type Service struct{}

func (self *Service) Add(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
		req  struct {
			id          string
			loadBalance string
		}
		service      global.ServiceType
		serviceBytes []byte
	)
	if err = filter.MSet(
		filter.El(&req.id, filter.FromString(ctx.Post("id"), "id").Required()),
		filter.El(&req.loadBalance, filter.FromString(ctx.Post("load_balance"), "load_balance")),
	); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	if _, exists := global.Services.Load(req.id); exists {
		resp["error"] = "服务ID已存在"
		return JSON(ctx, 400, &resp)
	}
	if req.loadBalance == "" {
		resp["error"] = "load_balance参数不能同时为空"
		return JSON(ctx, 400, &resp)
	}

	service.ID = req.id

	if serviceBytes, err = service.MarshalJSON(); err != nil {
		log.Err(err).Caller().Msg("将服务信息序列化成JSON字符串失败")
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}

	if err = global.Storage.SaveService(req.id, global.BytesToStr(serviceBytes)); err != nil {
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
			id          string
			loadBalance string
		}
		service      global.ServiceType
		serviceBytes []byte
	)
	if err = filter.MSet(
		filter.El(&req.id, filter.FromString(ctx.PathParams.Value("id"), "id").Required().Base64RawURLDecode()),
		filter.El(&req.loadBalance, filter.FromString(ctx.Post("load_balance"), "load_balance")),
	); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}

	if req.loadBalance == "" {
		resp["error"] = "load_balance参数不能同时为空"
		return JSON(ctx, 400, &resp)
	}

	service.ID = req.id

	if serviceBytes, err = service.MarshalJSON(); err != nil {
		log.Err(err).Caller().Msg("将服务信息序列化成JSON字符串失败")
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}

	if err = global.Storage.SaveService(req.id, global.BytesToStr(serviceBytes)); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Service) Delete(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
		id   string
	)
	if id, err = global.DecodeKey(ctx.PathParams.Value("id")); err != nil {
		return Status(ctx, 404)
	}
	if _, exist := global.Services.Load(id); !exist {
		return Status(ctx, 404)
	}
	err = global.Storage.DeleteStorageService(ctx.PathParams.Value("id"))
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
