package api

import "github.com/dxvgef/tsing"

// 设置路由
func SetRouter(engine *tsing.Engine) {
	// 检查secret
	router := engine.Group("", checkSecretFromHeader)

	var proxyHandler Engine
	router.GET("/engine/", proxyHandler.OutputJSON)
	router.POST("/engine/", proxyHandler.LoadAll)
	router.PUT("/engine/", proxyHandler.SaveAll)
}
