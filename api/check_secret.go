package api

import (
	"github.com/dxvgef/tsing"

	"local/global"
)

func checkSecretFromHeader(ctx *tsing.Context) error {
	if ctx.Request.Header.Get("SECRET") != global.Config.API.Secret {
		ctx.Abort()
		return Status(ctx, 401)
	}
	return nil
}
