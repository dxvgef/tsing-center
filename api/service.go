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

// 选取节点
func (self *Service) Select(ctx *tsing.Context) error {
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
	total := ci.Total()
	loseNodes := map[string]global.Node{}
	for i := 0; i < total; i++ {
		ip, port, expires := ci.Select()
		if expires > 0 && expires <= time.Now().Unix() {
			loseNodes[serviceID] = global.Node{
				IP:   ip,
				Port: port,
			}
			continue
		}
		resp["ip"] = ip
		resp["port"] = port
		break
	}
	if resp["ip"] == "" {
		return Status(ctx, http.StatusNotImplemented)
	}

	// 如果发现有过期的节点，则执行清理操作
	if len(loseNodes) > 0 {
		// 由于下面还有其它的逻辑，所以此处新开协程执行异步清理
		go global.Storage.Clean(loseNodes)
	}

	return JSON(ctx, 200, &resp)
}
