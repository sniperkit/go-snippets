
使用logrus库记录日志，增加了输出到kafka的hook。
 增加线上环境和本地环境的不同日志模式，在conf/conf.go文件中进行配置。
* 开发模式直接在控制台打印所有日志，会显示文件名和行号。支持颜色。
* 线上模式输出Info及以上的日志，输出到kafka


用法：
1. 将logrus.go拷贝到项目的log目录下。
1. 将conf/conf.go 和 app.conf拷贝到项目目录下，根据自己的情况修改app.conf
1. 调用log.Logger().Info("msg")来记录日志，支持Debug, Info, Warning, Error, Panic
