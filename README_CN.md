z这个项目汇集了用以演示[rpcx](https://github.com/smallnest/rpcx)的特性的例子。

每个例子都很简短，包括服务器代码和客户端代码，你可以在单机上进行测试。

- alias: 为服务起个别名
- authorization: 身份认证和授权
- clientpool: 客户端可以使用连接池，每个客户端都可以重用
- codec: 可以配置特定的编解码
- compression: 消息传递提供压缩解压缩的支持
- context: 服务器可以通过context获得客户端连接，比如客户端IP
- docker: 打包docker
- end2end: 点对点的代码，不经过注册中心
- etcd_registry: etcd注册中心的演示，使用etcd v2协议
- etcdv3_registry: etcd注册中心的演示，使用etcd v3协议
- geo: 根据经纬度坐标选择地理位置最近的数据中心
- group: 对服务进行分组，客户端只能访问特定组的服务
- java: JSON-RPC2的演示。服务器由rpcx提供，客户端由java实现
- kcp: 支持[kcp](https://github.com/skywind3000/kcp)
- multi_server: 多个服务器的负载均衡，无须注册中心
- nil_panic: 演示空指针错误，目前rpcx不能自动捕获空指针
- realservice: 我在colobu.com提供了一个rpcx的服务，你可以编写client进行测试。[弃用]
- reconnect: 服务在全宕机的情况客户端的重试功能
- reuseport: 在高版本内核的Linux操作系统上，你可以启动多个rpcx服务器，它们可以共用一个端口，这样一个实例宕机其它实例可以继续服务，Linux内核自己负载均衡
- rpcx-ui-demoserver: 演示rpcx管理界面
- timeout: 超时设置，客户端和服务器端都可以设置超时
- tls: 演示加密传输
- zookeeper_registry: zookeeper注册中心