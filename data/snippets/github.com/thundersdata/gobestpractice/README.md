# gobestpractice
golang 项目最佳实践。提供基本问题的建议写法，让大家能够专注于业务代码的实现。

编码规范参见 [Golang编码规范.md](https://github.com/thundersdata/gobestpractice/blob/master/Golang%E7%BC%96%E7%A0%81%E8%A7%84%E8%8C%83.md)

## Finished

* 配置文件读取(conf, 配置文件默认是项目根目录下的app.conf，yaml格式)
-- github.com/sirupsen/logrus
* 日志处理(log，依赖conf)
-- github.com/jinzhu/configor
* JSON Parser(parser/json)
-- github.com/tidwall/gjson
* YAML Parser(parser/yaml)
-- gopkg.in/yaml.v2

## TODO

1. Parser

   1. XML

1. DATABASE

   1. ORM
   1. HBase
   1. Redis

1. RPC

1. HTTP

1. Test & Mock

1. MessageQueue

   ​