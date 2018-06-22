## 多个基于alpine的SSR镜像

- http \
对ssr端口做了redirect转发 \
通常情况下，连接服务端口将被ssr转发至内部的nginx服务，显示一个静态网页。 \
客户端连接时，参数正确亦可正常使用。
- kcptun \
加装了kcptun的ssr，建议只对udp端口做映射，大多数配置可通过环境变量调整，详情见Dockerfile
- lite \
普通ssr，无任何特殊加成
