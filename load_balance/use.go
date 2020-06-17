package load_balance

import (
	"errors"
	"strings"

	"github.com/dxvgef/tsing-center/global"
	"github.com/dxvgef/tsing-center/load_balance/swrr"
	"github.com/dxvgef/tsing-center/load_balance/wr"
	"github.com/dxvgef/tsing-center/load_balance/wrr"
)

// 使用指定算法构建负载均衡器的实例
func Build(name string) (global.LoadBalance, error) {
	name = strings.ToUpper(name)
	switch name {
	// 加权随机
	case "WR":
		return wr.New(), nil
	// 加权轮循
	case "WRR":
		return wrr.New(), nil
	// 平滑加权轮循
	case "SWRR":
		return swrr.New(), nil
	}
	return nil, errors.New("不支持的负载均衡算法")
}
