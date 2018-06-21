说到rpc让我想起了刚毕业面试的时候，被问到是否了解rpc？我记得当时我的回答是“课本上学过rpc，只知道是远程过程调用，但没有用过，具体也不知道是什么”。的确，大学中间件这门课程里有讲到rpc，里面还引入了一个非常难理解的概念——“桩”，英文应该叫”stub”。现在的rpc实现里，stub这个概念好像都没见到了，应该都是叫”method”。

实现一个rpc服务器很难吗？rpc服务器也就是在tcp服务器的基础上加上自定义的rpc协议而已。一个rpc协议里，主要有个3个非常重要的信息。

调用的远程method名字，一般就是一个函数名
call参数，也就是发送给服务器的数据
客户端生成的调用请求seq
除了最后一点，其他两点显然就是组成一个普通的函数调用而已，这也就是远程过程调用了。最后一点的seq只是rpc内部实现的一个关键点而已，也许还有其他的实现方式，而不是依赖seq来保证单连接的并发请求。

####如何用C语言实现rpc服务器？ 用C语言来实现rpc服务器，先假设我们从socket里读取到了一个来自于客户端的call请求，这个call请求里面封装了上面提到的3点信息。

执行这个call，最重要的就是——“从call里取出method，也就是一个函数名(字符串)，然后要通过这个函数名去执行对应的函数”。

C语言由于没有反射机制，于是不能通过函数的名字去调对应的函数；因此，可以使用一张hash表保存所有远程函数的名字到函数的对应关系。这样就可以通过method查找一次hash表得到真正的函数，接下来就可以执行函数了，函数的执行结果当然就是作为响应返回给客户端。

当然，这些需要被调用执行的函数都是在服务器初始化的时候事先注册到这张hash表中的。到这里，感觉一切都结束了。其实，还有参数(请求数据)和返回值(应答数据)，这部分主要是涉及到序列化算法，留到下次的“反射和序列化”再聊了。

####使用Go RPC框架写一个简单服务器

定义服务

type EchoServer bool

func (s *EchoServer) Echo(req *string, res *string) error {
        res = req
                return nil
}
注册服务

echo := new(EchoServer)
    if err := rpc.Register(echo); err != nil {
            return err
    }

// 下面基本是固定代码，创建一个tcp服务器
var listener *net.TCPListener
if tcpAddr, err := net.ResolveTCPAddr("tcp", addr); err != nil {
        return err
} else {
        if listener, err = net.ListenTCP("tcp", tcpAddr); err != nil {
            return err
        }
}

for {
        conn, err := listener.Accept()
                if err != nil {
                    sc.LOG.Error(err)
                        continue
                }
        // 服务器走rpc框架的入口，也就是在一个tcp连接上采用rpc协议来处理请求和应答。
        go rpc.ServeConn(conn)
}
return nil
这样，就实现了一个简单的echo服务器，功能逻辑都在自定义的EchoServer里面实现。接下来只需要将EchoServer的一个实例对象注册到rpc框架中，客户端就可以直接调用EchoServer对象中的方法了。

Go rpc服务器的内部实现在思路上和前文提到的C语言实现基本是一样的。细心的人可能注意到了，这里注册的是EchoServer对象，而不是EchoSever对象的方法，然后前文的C语言实现中注册的直接是函数，那么rpc服务器如何能够根据method去调用对应的方法呢？Go语言在这里其实采用反射的手段，虽然表面上是注册的EchoServer对象，实际却是通过反射取得了EchoServer的所有方法，然后采用了map表保存了method到方法的映射，这就回到了前文C语言的实现思路中去了。如果没有反射的支持，就只能一个一个的方法全部注册一遍了，并且代码组织上也不够优雅。

####编码解码器 rpc客户端有一个编码解码器定义了如何发送请求和读取应答，那么服务器端必然有一个编码解码器定义了如何读取请求和发送应答，刚好是一个相反的过程，这也就是序列化和反序列化的一个用处。

type ServerCodec interface {
    ReadRequestHeader(*Request) error
        ReadRequestBody(interface{}) error
        WriteResponse(*Response, interface{}) error

        Close() error
} 和客户端一样，你只需要实现ServerCodec这个接口，就可以自定义服务端的编码解码器，实现自定义的rpc协议了。 Go rpc服务器端和客户端都是默认使用的Gob序列化协议数据。Go里面为什么采用自实现的Gob，而不是Google的protobuf，这个也留到下次的*"反射和序列化"*聊吧。
注意，ServerCodec接口中的ReadRequestBody(interface{}) error方法主要是用来读取客户端call请求的参数数据，也就是将socket中读出来的数据包解码到interface{}所指向的具体对象中去，这里的interface{}可以理解为C语言中的void *。你一定注意到了，不同的rpc服务定义的参数类型完全不同，并且在rpc框架内部都采用了interface{}来适配，那么框架内部如何知道读取的socket数据要解码到什么具体类型中去呢？这里又是涉及到了反射，因为有了反射就可以从interface{}得到具体的类型。

接下来真得好好的说说“反射和序列化/反序列化”了。

####几个内部细节 * 服务器在发送应答的时候，同样采用了一把锁来保证原子性写入socket * 请求/应答结构(ServerCodec接口中出现的Request/Response)采用链表实现一个pool，需要一个Requet/Response的时候都从free list中获取，不再频繁的分配。这一点只是一个优化。 * Go的rpc除了直接跑在tcp服务器上，还可以跑在更高层一点的http服务器上。 * Go的rpc调用也提供了json序列化。


准备改变一下以前写技术文章老是大段大段代码的风格，试试小代码能不能将自己的想法交代清楚，哈哈。这一篇就先到此为止了。

下一篇就看一下Go RPC框架的性能，做一个benchmark比较下。
