# tsing-center

使用Go开发的服务中心，提供服务注册、发现、负载均衡、熔断等服务治理功能，并且可以方便的横向扩展服务节点实现自身的负载均衡。

## 应用场景

在分布式架构中，服务中心通常用于各架构内部各服务之间的互相发现，例如API网关和业务服务节点之间，使API网关能正确的将客户端请求反向代理到具体的某个服务节点进程。

业务逻辑的流程如下：
1. 客户端请求到达API网关
2. API网关根据路由匹配到了服务
3. API网关根据服务的设置向`Tsing Center`发送获取目标节点的请求
4. `Tsing Center`通过负载均衡策略选取出一个目标节点，并返回给API网关
5. API网关将客户端请求反向代理到目标节点


## 安装方法

至Release页面下载最新版本的压缩包，解压后编辑默认的`config.yml`文件，并运行`tsing-center`或`tsing-center.exe`二进制文件。
