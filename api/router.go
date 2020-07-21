package api

import "github.com/dxvgef/tsing"

// 设置路由
func SetRouter(engine *tsing.Engine) {
	// 检查secret中间件
	router := engine.Group("", checkSecretFromHeader)

	router.GET("/ip", GetIP) // 用于客户端获取IP地址

	// 数据管理
	var dataHandler Data
	router.GET("/data/", dataHandler.OutputJSON) // 将本节点所有本地缓存数据以JSON格式输出
	router.POST("/data/", dataHandler.LoadAll)   // 从存储器加载所有数据到本地缓存
	router.PUT("/data/", dataHandler.SaveAll)    // 将本节点所有本地缓存数据写入到存储器

	// 服务管理
	var serviceHandler Service
	router.POST("/services/", serviceHandler.Add)                    // 创建服务
	router.PUT("/services/:serviceID", serviceHandler.Put)           // 重写或创建服务
	router.GET("/services/:serviceID/select", serviceHandler.Select) // 获取服务中的节点信息
	router.DELETE("/services/:serviceID", serviceHandler.Delete)     // 删除服务

	// 节点管理
	var nodeHandler Node
	router.POST("/nodes/", nodeHandler.Add)                           // 创建节点
	router.PUT("/nodes/:serviceID/:node", nodeHandler.Put)            // 重写或创建节点
	router.DELETE("/nodes/:serviceID/:node", nodeHandler.Delete)      // 删除节点
	router.PATCH("/nodes/:serviceID/:node/:attrs", nodeHandler.Patch) // 更新节点属性
	router.POST("/nodes/:serviceID/:node", nodeHandler.Touch)         // 节点触活
}
