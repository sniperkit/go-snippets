# 基础规范
[TOC]

## gofmt(goimports)

大部分的格式问题可以通过gofmt解决，gofmt自动格式化代码，保证所有的go代码与官方推荐的格式保持一致，于是所有格式有关问题，都以gofmt的结果为准。
## 项目目录结构
```
bestpractice 包名，用全小写。单词之间直接连接，不需要用下划线等
├── README.md 介绍软件及文档入口
├── doc 该项目的文档
├── main.go 项目主文件
├── xxx.go 项目文件
├── .gitignore
├── Gopkg.lock 依赖管理文件
├── Gopkg.toml 依赖管理文件
├── run.sh 可选，项目执行文件
├── vendor 存放第三方库和公司公共库
└── model 项目子目录结构
    ├── db.go
    └── user.go
```
## 包管理器
统一使用[dep](https://golang.github.io/dep/docs/introduction.html)进行包管理，需要**科学上网**。
* 安装： ` go get -u github.com/golang/dep/cmd/dep `
* 初始化： ` dep init `
* 查看依赖状态： ` dep status `
* 添加依赖： `dep ensure -add xxx/xxx `
* 更新依赖： ` dep ensure `
## 注释
可以通过 /* …… */ 或者 // ……增加注释， //之后应该加一个空格。

在编码阶段应该同步写好变量、函数、包的注释，最后可以利用godoc导出文档。注释必须是完整的句子，句子的结尾应该用句号作为结尾（英文句号）。
每个包都应该有一个包注释，一个位于package子句之前的块注释或行注释。包如果有多个go文件，只需要出现在一个go文件中即可。
```go
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package net provides a portable interface for network I/O, including
TCP/IP, UDP, domain name resolution, and Unix domain sockets.
......
*/
package net
......
```
导出函数注释，第一条语句应该为一条概括语句，并且使用被声明的名字作为开头。
```go
// Myfunction return the sum of a and b.
func Myfunction(sum int) (a, b int) {
```
## 命名
> *Finally, the convention in Go is to use MixedCaps or mixedCaps rather than underscores to write multiword names.*

Go倾向用驼峰是命名法，而不是用下划线分割单词。

###文件名

全小写，单词之间用下划线相连。

### 包名

> *By convention, packages are given lower case, single-word names; there should be no need for underscores or mixedCaps.*

全小写。单词之间直接连接，不需要用下划线或者单词首字母大写。
### 接口名
> *By convention, one-method interfaces are named by the method name plus an -er suffix or similar modification to construct an agent noun: Reader, Writer, Formatter, CloseNotifier etc.*

单个函数的接口名以”er”作为后缀，如Reader,Writer。接口的实现则去掉“er”。
```go
type Reader interface {
        Read(p []byte) (n int, err error)
}
```
两个函数的接口名综合两个函数名
```go
type WriteFlusher interface {
    Write([]byte) (int, error)
    Flush() error
}
```
三个以上函数的接口名，类似于结构体名
```go
type Car interface {
    Start([]byte)
    Stop() error
    Recover()
}
```
### 函数名
驼峰式。
### 变量
驼峰式。局部变量用小写字母开头。需要在package外部使用的全局变量用大写字母开头，否则用小写字母开头。
### 常量
采用驼峰法 `maxLength` 或者 `MaxLength`。 不要用下划线。
## 错误处理
* error作为函数的值返回,必须对error进行处理
* 错误描述如果是英文必须为小写，不需要标点结尾
* 采用独立的错误流进行处理
    * 不要采用这种方式

    ```go
    if err != nil {
        // error handling
    } else {
        // normal code
    }
    ```
    * 而要采用下面的方式

    ```go
    if err != nil {
        // error handling
        return // or continue, etc.
    }
    // normal code
    ```
在逻辑处理中禁用panic。在main包中只有当实在不可运行的情况采用panic，例如文件无法打开，数据库无法连接导致程序无法 正常运行，但是对于其他的package对外的接口不能有panic，只能在包内采用。 建议在main包中使用log.Fatal来记录错误，这样就可以由log来结束程序。
## 单元测试
* 单元测试文件名命名规范为 example_test.go
* 测试用例的函数名称必须以 Test 开头，例如：TestExample

---

# 进阶规范
### 在项目中不要使用相对路径引入包，要用绝对路径
    ​```go
    // 这是不好的导入
    import "../net"
    // 这是正确的做法
    import "xxx.com/proj/net"
    ​```
### 声明slice
声明空的slice应该使用下面的格式:
```go
var t []string
```
而不是这种格式:
```go
t := []string{}
```
前者声明了一个nil slice而后者是一个长度为0的非nil的slice。
### 非空slice检查
不要使用下面的方式检查空的slice:
```go
if s != nil && len(s) > 0 {
    ...
}
```
直接比较长度即可：
```go
if len(s) > 0 {
    ...
}
```
同样的道理也适用 map和channel。
### 关于字符串大小写
错误字符串不应该大写。
应该写成：
```go
fmt.Errorf("failed to write data")
```
而不是写成：
```go
fmt.Errorf("Failed to write data")
```
这是因为这些字符串可能和其它字符串相连接，组合后的字符串如果中间有大写字母开头的单词很突兀，除非这些首字母大写单词是固定使用的单词。

缩写词必须保持一致，比如都大写URL或者小写url。比如HTTP、ID等。
### 方法接收器
* 名称一般采用struct的第一个字母，且为小写，而不是this，me或者self。
    ```go
    type T struct{}
    func (t *T)Get(){}
    ```
* 如果接收者是map,slice或者chan，不要用指针传递。如果需要对slice进行修改，通过返回值的方式重新赋值。
* 如果接收者是含有sync.Mutex或者类似同步字段的结构体，必须使用指针传递避免复制。
* 如果接收者是大的结构体或者数组，使用指针传递会更有效率。
### package级的Error变量
通常会把自定义的Error放在package级别中，统一进行维护:
```go
var (
    ErrCacheMiss = errors.New("memcache: cache miss")
    ErrCASConflict = errors.New("memcache: compare-and-swap conflict")
    ErrNotStored = errors.New("memcache: item not stored")
    ErrServerError = errors.New("memcache: server error")
    ErrNoStats = errors.New("memcache: no statistics available")
    ErrMalformedKey = errors.New("malformed: key is too long or contains invalid characters")
    ErrNoServers = errors.New("memcache: no servers configured or available")
)
```
并且变量以**Err**开头。
### 空字符串检查
不要使用下面的方式检查空字符串:
```go
if len(s) == 0 {
	...
}
```
而是使用下面的方式
```go
if s == "" {
    ...
}
```
### 省略不必要的变量
比如
```go
var whitespaceRegex, _ = regexp.Compile("\\s+")
```
可以简写为
```go
var whitespaceRegex = regexp.MustCompile(`\s+`)
```

---

# 工具
## GoLand
[官方网站](https://www.jetbrains.com/go/)
激活参考: http://blog.csdn.net/lz_1990/article/details/79360102

## gometalinter
[官方网站](https://github.com/alecthomas/gometalinter)

goland配置metalinter:
1. 安装metalinter: ``` go get -u github.com/alecthomas/gometalinter ```
1. 安装所有的linter: ``` gometalinter --install ```
1. 配置goland。(也可以直接执行 ``` gometalinter <filename> ```)

![1.png](https://upload-images.jianshu.io/upload_images/10018337-a2814c92056bca53.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

![2.png](https://upload-images.jianshu.io/upload_images/10018337-f86ebb16a1c85ae9.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)
## godoc
[官方网站](https://godoc.org/golang.org/x/tools/cmd/godoc)
* 安装: ``` go get golang.org/x/tools/cmd/godoc ```
* 使用： ``` godoc -server=:6060 ``` 或者 ``` godoc <packagename> ```
# 参考资料
* [http://golang.org/doc/effective_go.html](http://golang.org/doc/effective_go.html)
* [https://code.google.com/p/go-wiki/wiki/CodeReviewComments](https://code.google.com/p/go-wiki/wiki/CodeReviewComments)
* https://blog.csdn.net/myzlhh/article/details/52269591
* https://sheepbao.github.io/post/golang_code_specification/
* http://colobu.com/2017/02/07/write-idiomatic-golang-codes/
