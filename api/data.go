package api

import (
	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-center/global"
)

type Data struct {
	Services []global.ServiceType         `json:"services,omitempty"`
	Nodes    map[string][]global.NodeType `json:"nodes,omitempty"`
}

func (self *Data) OutputJSON(ctx *tsing.Context) error {
	var data Data
	global.Services.Range(func(_, value interface{}) bool {
		v, ok := value.(global.ServiceType)
		if !ok {
			log.Error().Caller().Msg("类型断言失败")
			return false
		}
		data.Services = append(data.Services, global.ServiceType{
			ID:          v.ID,
			LoadBalance: v.LoadBalance,
		})
		return true
	})
	if data.Nodes == nil {
		data.Nodes = make(map[string][]global.NodeType, global.SyncMapLen(&global.Nodes))
	}
	global.Nodes.Range(func(key, value interface{}) bool {
		serviceID := key.(string)
		lb, ok := value.(global.LoadBalance)
		if !ok {
			log.Error().Caller().Msg("类型断言失败")
			return false
		}
		data.Nodes[serviceID] = lb.Nodes()
		return true
	})
	bs, err := data.MarshalJSON()
	if err != nil {
		ctx.ResponseWriter.WriteHeader(500)
		return ctx.Caller(err)
	}
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ctx.ResponseWriter.WriteHeader(200)
	if _, err = ctx.ResponseWriter.Write(bs); err != nil {
		log.Err(err).Caller().Send()
		return ctx.Caller(err)
	}
	return nil
}

func (*Data) LoadAll(ctx *tsing.Context) error {
	resp := make(map[string]string)
	if err := loadAll(); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
func (*Data) SaveAll(ctx *tsing.Context) error {
	resp := make(map[string]string)
	if err := saveAll(); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

// 加载所有数据
func loadAll() (err error) {
	return global.Storage.LoadAll()
}

// 保存所有数据
func saveAll() (err error) {
	return global.Storage.SaveAll()
}
