## 多个基于alpine的v2ray镜像

- Lite \
普通v2ray，无任何特殊加成，可通过环境变量修改端口和UUID
- websocket（待更新） \
底层传输协议使用websocket，利用nginx做了内部转发 \
设置正确的二级目录即可正常使用。
- mKcp（待更新） \
加装了kcptun的v2ray
- http（待更新）
底层传输协议使用HTTP/2.0
- heroku \
可用于heroku平台，支持tls加密
