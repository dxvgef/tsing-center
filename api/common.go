package api

import (
	"strings"

	"github.com/dxvgef/tsing"
)

// 用于客户端获取IP地址
func GetIP(ctx *tsing.Context) error {
	return String(ctx, 200, ctx.Request.RemoteAddr[:strings.Index(ctx.Request.RemoteAddr, ":")])
}
