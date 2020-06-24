# tsing-center

`Tsing Center`是一个开源、跨平台、去中心化集群、动态配置的服务中心。

## 应用场景

在分布式架构中，服务中心通常用于架构内部的网络服务进程之间的互相发现(IP:端口)。

例如API网关和业务服务节点之间，可以通过服务中心来互相发现，使API网关能正确的将客户端请求反向代理到最终的上游服务节点。

## 功能特性
- 服务注册，通过API动态注册服务的节点信息
- 服务发现，通过API或DNS查询获取服务的节点信息
- 负载均衡，对服务中的节点使用负载均衡算法进行选取
- 健康检查，通过API刷新节点的生命周期(心跳)，自动剔除"心跳"超时的节点
- 去中心化集群，轻松组建横向扩展的服务中心集群，并用任意节点做请求入口
- API动态配置，可通过RESTful和gRPC协议的API对配置进行动态变更，无需重启进程
- 持久存储，支持`etcd`、`consul`、`redis`多种数据源

### 存储引擎
- [x] etcd
- [ ] consul
- [ ] redis

### 负载均衡
- [x] SWRR，平滑加权轮循，类似Nginx
- [x] WRR，加权轮循，类似LVS
- [x] WR，加权随机

## 相关资源

- [Tsing](https://github.com/dxvgef/tsing) 高性能、微核心的Go语言HTTP服务框架
- [Tsing Gateway](https://github.com/dxvgef/tsing-gateway) 开源、跨平台、去中心化集群、动态配置的API网关

## 用户及案例

如果你在使用本项目，请通过[Issues](https://github.com/dxvgef/tsing-center/issues)告知我们项目的简介

## 帮助/说明

本项目处于开发初期阶段，API和数据存储结构可能会频繁变更，暂不建议在生产环境中使用，如有问题可在[Issues](https://github.com/dxvgef/tsing-center/issues)里提出。

诚邀更多的开发者为本项目开发管理面板和官方网站等资源，帮助这个开源项目更好的发展。
