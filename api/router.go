package api

import "github.com/dxvgef/tsing"

// 设置路由
func SetRouter(engine *tsing.Engine) {
	// 检查secret
	router := engine.Group("", checkSecretFromHeader)

	// 数据管理
	var dataHandler Data
	router.GET("/data/", dataHandler.OutputJSON)
	router.POST("/data/", dataHandler.LoadAll)
	router.PUT("/data/", dataHandler.SaveAll)

	// 服务管理
	var serviceHandler Service
	router.POST("/service/", serviceHandler.Add)
	router.PUT("/service/:serviceID", serviceHandler.Put)
	router.DELETE("/service/:serviceID", serviceHandler.Delete)

	// 节点管理
	var nodeHandler Service
	router.POST("/node/", nodeHandler.Add)
	router.PUT("/node/:serviceID/:node", nodeHandler.Put)
	router.DELETE("/node/:serviceID/:node", nodeHandler.Delete)
}
