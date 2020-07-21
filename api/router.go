package api

import "github.com/dxvgef/tsing"

// 设置路由
func SetRouter(engine *tsing.Engine) {
	// 检查secret
	router := engine.Group("", checkSecretFromHeader)

	// 用于客户端获取IP地址
	router.GET("/ip", GetIP)

	// 数据管理
	var dataHandler Data
	router.GET("/data/", dataHandler.OutputJSON)
	router.POST("/data/", dataHandler.LoadAll)
	router.PUT("/data/", dataHandler.SaveAll)

	// 服务管理
	var serviceHandler Service
	router.POST("/services/", serviceHandler.Add)
	router.PUT("/services/:serviceID", serviceHandler.Put)
	router.GET("/services/:serviceID/select", serviceHandler.Select)
	router.DELETE("/services/:serviceID", serviceHandler.Delete)

	// 节点管理
	var nodeHandler Node
	router.POST("/nodes/", nodeHandler.Add)
	router.PUT("/nodes/:serviceID/:node", nodeHandler.Put)
	router.DELETE("/nodes/:serviceID/:node", nodeHandler.Delete)
	router.PATCH("/nodes/:serviceID/:node/:attrs", nodeHandler.Patch)
}
