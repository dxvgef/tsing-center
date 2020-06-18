package api

import (
	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-center/global"
)

type Data struct {
	Services []global.ServiceConfig   `json:"services,omitempty"`
	Nodes    map[string][]global.Node `json:"nodes,omitempty"`
}

func (self *Data) OutputJSON(ctx *tsing.Context) error {
	var data Data
	global.Services.Range(func(_, value interface{}) bool {
		v, ok := value.(global.Cluster)
		if !ok {
			log.Error().Caller().Msg("类型断言失败")
			return false
		}
		config := v.Config()
		data.Services = append(data.Services, config)
		if data.Nodes == nil {
			data.Nodes = map[string][]global.Node{}
		}
		nodes := v.Nodes()
		for k := range nodes {
			data.Nodes[config.ServiceID] = append(data.Nodes[config.ServiceID], nodes[k])
		}
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
