package cluster

import (
	"errors"
	"strings"

	"github.com/dxvgef/tsing-center/cluster/swrr"
	"github.com/dxvgef/tsing-center/cluster/wr"
	"github.com/dxvgef/tsing-center/cluster/wrr"
	"github.com/dxvgef/tsing-center/global"
)

// 构建一个使用指定负载均衡算法的集群实例
func Build(config global.ServiceConfig) (global.Cluster, error) {
	config.LoadBalance = strings.ToUpper(config.LoadBalance)
	switch config.LoadBalance {
	// 加权随机
	case "WR":
		return wr.New(config), nil
	// 加权轮循
	case "WRR":
		return wrr.New(config), nil
	// 平滑加权轮循
	case "SWRR":
		return swrr.New(config), nil
	}
	return nil, errors.New("不支持的集群负载均衡规则")
}
