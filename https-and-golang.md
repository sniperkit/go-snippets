https原理以及golang基本实现 (http://studygolang.com/articles/4461)
 关于https
 背景知识
 密码学的一些基本知识
 大致上分为两类，基于key的加密算法与不基于key的加密算法。现在的算法基本都是基于key的，key就以一串随机数数，更换了key之后，算法还可以继续使用。

 基于key的加密算法又分为两类，对称加密和不对称加密，比如DES,AES那种的，通信双方一方用key加密之后，另一方用相同的key进行反向的运算就可以解密。

 不对称加密比较著名的就是RSA,加密的时候有一个公钥和一个私钥，公钥是可以交给对方的，a给b发送信息，a用自己的私钥加密，b用a的公钥解密，反之，b给a发送信息，b用自己的私钥加密。

 在通信之前，需要经过一些握手的过程，双方交换公钥，这个就是key exchange的过程，https最开始的阶段就包含了这个key exchange的过程，大概原理是这样，有些地方还要稍微复杂一些。

 数字证书与CA
 数字证书相当于是服务器的一个“身份证”，用于唯一标识一个服务器。一般而言，数字证书从受信的权威证书授权机构 (Certification Authority，证书授权机构)买来的（免费的很少），浏览器里面一般就内置好了一些权威的CA，在使用https的时候，只要是这些CA签发的证书，浏览器都是可以认证的，要是在与服务器通信的时候，收到一个没有权威CA认证的证书，就会报出提醒不受信任证书的错误，就像登录12306一样，但是也可以选择接受。

 在自己的一些项目中，通常是自己签发一个ca根证书，之后这个根证书签发一个server.crt，以及server.key给服务端，server.key是服务端的私钥，server.crt包含了服务端的公钥还有服务端的一些身份信息。在客户端和服务端通信的时候（特别是使用代码编写的客户端访问的时候），要指定ca根证书，作用就相当于是浏览器中内置的那些权威证书一样，用于进行服务端的身份检测。

 证书的格式：

 ca证书在为server.crt证书签名时候的大致流程参考这个(http://www.tuicool.com/articles/aymYbmM)：

 数字证书由两部分组成：

 1、C：证书相关信息（对象名称+过期时间+证书发布者+证书签名算法….）

 2、S：证书的数字签名 （由CA证书通过加密算法生成的）

 其中的数字签名是通过公式S = F(Digest(C))得到的。

 Digest为摘要函数，也就是 md5、sha-1或sha256等单向散列算法，用于将无限输入值转换为一个有限长度的“浓缩”输出值。比如我们常用md5值来验证下载的大文件是否完整。大文件的内容就是一个无限输入。大文件被放在网站上用于下载时，网站会对大文件做一次md5计算，得出一个128bit的值作为大文件的摘要一同放在网站上。用户在下载文件后，对下载后的文件再进行一次本地的md5计算，用得出的值与网站上的md5值进行比较，如果一致，则大 文件下载完好，否则下载过程大文件内容有损坏或源文件被篡改。这里还有一个小技巧常常在机器之间copy或者下载压缩文件的时候也可以用md5sum的命令来进行检验，看看文件是否完整。

 F为签名函数。CA自己的私钥是唯一标识CA签名的，因此CA用于生成数字证书的签名函数一定要以自己的私钥作为一个输入参数。在RSA加密系统中，发送端的解密函数就是一个以私钥作为参数的函数，因此常常被用作签名函数使用。因此CA用私钥解密函数作为F，以CA证书中的私钥进行加密，生成最后的数字签名，正如最后一部分实践时候给出的证书生成过程，生成server.crt的时候需要ca.crt（包含根证书的信息）和ca.key（根证书的私钥）都加入进去。

 接收端接收服务端数字证书后，如何验证数字证书上携带的签名是这个CA的签名呢？当然接收端首先需要指定对应的CA，接收端会运用下面算法对数字证书的签名进行校验：
 F'(S) ?= Digest(C)

 接收端进行两个计算，并将计算结果进行比对：

 1、首先通过Digest(C)，接收端计算出证书内容（除签名之外）的摘要，C的内容都是明文可以看到到的。

 2、数字证书携带的签名是CA通过CA密钥加密摘要后的结果，因此接收端通过一个解密函数F'对S进行“解密”。就像最开始介绍的那样，在RSA系统中，接收端使用CA公钥（包含在ca.crt中）对S进行“解密”，这恰是CA用私钥对S进行“加密”的逆过程。

 将上述两个运算的结果进行比较，如果一致，说明签名的确属于该CA，该证书有效，否则要么证书不是该CA的，要么就是中途被人篡改了。

 对于self-signed(自签发)证书来说，接收端并没有你这个self-CA的数字证书，也就是没有CA公钥，也就没有办法对数字证书的签名进行验证。因此如果要编写一个可以对self-signed证书进行校验的接收端程序的话，首先我们要做的就是建立一个属于自己的CA，用该CA签发我们的server端证书，之后给客户端发送信息的话，需要对这个根证书进行指定，之后按上面的方式进行验证。

 可以使用openssl x509 -text -in client.crt -noout 查看某个证书文件所包含的具体信息。

 HTTPS基本过程概述

 https协议是在http协议的基础上组成的secure的协议。主要功能包含一下两个方面:

 1 通信双方的身份认证

 2 通信双方的通信过程加密

 下面通过详细分析https的通信过程来解释这两个功能。

 具体参考这两个文章：

 http://www.fenesky.com/blog/2014/07/19/how-https-works.html
http://www.ruanyifeng.com/blog/2014/02/ssl_tls.html

1、client 发送 sayhello给server端，说明client所支持的加密套件，还有一个随机数1。
2、server 发送 sayhello给client端，端把server.crt发送给客户端,server.crt采用还有一个随机数2。
3、client端生成preMaster key 这个是随机数3，之后三个随机数结合在一起生成MasterSecret,之后生成session secret，使用指定的ca进行身份认证，就像之前介绍的那样，都正常的话，就切换到加密模式。
4、client端使用server.crt中的公钥对preMasterSecret进行加密，如果要进行双向认证的话，client端会把client.crt一并发送过去，server端接受到数据，解密之后，也有了三个随机数，采用同样的方式，三个随机数生成通信所使用的session secret。具体session secret的结构可以参考前面列出的两个博客。server端完成相关工作之后，会发一个ChangeCipherSpec给client，通知client说明自己已经切换到相关的加解密模式，之后发一段加密信息给client看是否正常。
5、client端解密正常，之后就可以按照之前的协议，使用session secret进行加密的通信了。

整体看下，开始的时候建立握手的过程就是身份认证的过程，之后认证完毕之后，就是加密通信的过程了，https的两个主要做用就实现了。

相关实践

比较典型的证书生成的过程：

openssl genrsa -out ca.key 2048

#这里可以使用 -subj 不用进行交互 当然还可以添加更多的信息
openssl req -x509 -new -nodes -key ca.key -subj "/CN=zju.com" -days 5000 -out ca.crt

openssl genrsa -out server.key 2048

#这里的/cn可以是必须添加的 是服务端的域名 或者是etc/hosts中的ip别名
openssl req -new -key server.key -subj "/CN=server" -out server.csr

openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 5000
注意生成client端证书的时候，注意要多添加一个字段，golang的server端认证程序会对这个字段进行认证：

openssl genrsa -out client.key 2048

openssl req -new -key client.key -subj "/CN=client" -out client.csr

echo extendedKeyUsage=clientAuth > extfile.cnf

openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -extfile extfile.cnf -out client.crt -days 5000 
https客户端和服务端单向校验

这部分参考了这个（http://www.tuicool.com/articles/aymYbmM
），里面代码部分讲得比较细致。

服务端采用证书，客户端采用普通方式访问：

//server端代码
package main

import (
	"fmt"
	"net/http"
	"os"
       )

func handler(w http.ResponseWriter, r *http.Request) {
    	fmt.Fprintf(w,
		"Hi, This is an example of https service in golang!")
}

func main() {
    	http.HandleFunc("/", handler)
	    //http.ListenAndServe(":8080", nil)
	    _, err := os.Open("cert_server/server.crt")
	    if err != nil {
		panic(err)
	    }
	http.ListenAndServeTLS(":8081", "cert_server/server.crt",
		"cert_server/server.key", nil)
}
client端直接发请求，什么都不加，会报如下错误：

2015/07/11 18:13:50 http: TLS handshake error from 10.183.47.203:58042: remote error: bad certificate
使用浏览器直接访问的话，之后点击信赖证书，这个时候就可以正常get到消息

或者使用curl -k https:// 来经行访问，相当于忽略了第一步的身份验证的工作。
要是不加-k的话 使用curl -v 参数打印出来详细的信息，会看到如下的错误：

curl: (60) SSL certificate problem: Invalid certificate chain
说明是认证没有通过，因为客户端这面并没有提供可以信赖的根证书来对服务端发过来的证书进行验，/CN使用的直接是ip地址，就会报下面的错误：

Get https://10.183.47.206:8081: x509: cannot validate certificate for 10.183.47.206 because it doesn't contain any IP SANs
最好是生成证书的时候使用域名，或者是在/etc/hosts中加上对应的映射。

可以发送请求的客户端的代码如下，注意导入根证书的方式:

package main

import (
	//"io"
	//"log"
	"crypto/tls"
	"crypto/x509"
	//"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	//"strings"
       )

func main() {
    //x509.Certificate.
pool := x509.NewCertPool()
	  //caCertPath := "etcdcerts/ca.crt"
	  caCertPath := "certs/cert_server/ca.crt"

	  caCrt, err := ioutil.ReadFile(caCertPath)
	  if err != nil {
	      fmt.Println("ReadFile err:", err)
		  return
	  }
      pool.AppendCertsFromPEM(caCrt)
	  //pool.AddCert(caCrt)

	  tr := &http.Transport{
TLSClientConfig:    &tls.Config{RootCAs: pool},
			DisableCompression: true,
	  }
client := &http.Client{Transport: tr}

	resp, err := client.Get("https://server:8081")

	    if err != nil {
		panic(err)
	    }

	body, _ := ioutil.ReadAll(resp.Body)
	    fmt.Println(string(body))
	    fmt.Println(resp.Status)
}
使用curl命令的话，就加上--cacrt ca.crt证书，这样就相当于添加了可信赖的证书，身份认证的操作就可以成功了。

比如生成服务端证书的时候/CN写的是server 那client发送的时候也发送给https://server:8081就好，不过在本地的/etc/hosts中要加上对应的映射。

客户端和服务端的双向校验：

按照之前的方式，客户端生成证书，根证书就按之前的那个：

openssl genrsa -out client.key 2048

openssl req -new -key client.key -subj "/CN=client" -out client.csr

openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 5000
server端代码进行改进，添加受信任的根证书。

// gohttps/6-dual-verify-certs/server.go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
       )

type myhandler struct {
}

func (h *myhandler) ServeHTTP(w http.ResponseWriter,
	r *http.Request) {
    fmt.Fprintf(w,
	    "Hi, This is an example of http service in golang!\n")
}

func main() {
pool := x509.NewCertPool()
	  caCertPath := "cert_server/ca.crt"

	  caCrt, err := ioutil.ReadFile(caCertPath)
	  if err != nil {
	      fmt.Println("ReadFile err:", err)
		  return
	  }
      pool.AppendCertsFromPEM(caCrt)

	  s := &http.Server{
Addr:    ":8081",
	 Handler: &myhandler{},
	 TLSConfig: &tls.Config{
ClientCAs:  pool,
	    ClientAuth: tls.RequireAndVerifyClientCert,
	 },
	  }

      err = s.ListenAndServeTLS("cert_server/server.crt", "cert_server/server.key")
	  if err != nil {
	      fmt.Println("ListenAndServeTLS err:", err)
	  }
}
客户端代码改进，发送的时候把指定client端的client.crt以及client.key

// gohttps/6-dual-verify-certs/client.go

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
       )

func main() {
pool := x509.NewCertPool()
	  caCertPath := "certs/cert_server/ca.crt"

	  caCrt, err := ioutil.ReadFile(caCertPath)
	  if err != nil {
	      fmt.Println("ReadFile err:", err)
		  return
	  }
      pool.AppendCertsFromPEM(caCrt)

	  cliCrt, err := tls.LoadX509KeyPair("certs/cert_server/client.crt", "certs/cert_server/client.key")
	  if err != nil {
	      fmt.Println("Loadx509keypair err:", err)
		  return
	  }

tr := &http.Transport{
TLSClientConfig: &tls.Config{
RootCAs:      pool,
		  Certificates: []tls.Certificate{cliCrt},
		 },
    }
client := &http.Client{Transport: tr}
	resp, err := client.Get("https://server:8081")
	    if err != nil {
		fmt.Println("Get error:", err)
		    return
	    }
	defer resp.Body.Close()
	    body, err := ioutil.ReadAll(resp.Body)
	    fmt.Println(string(body))
}
但实际上，这样是不行的，server端会报这样的错误：


client's certificate's extended key usage doesn't permit it to be used for client authentication
因为client的证书生成方式有一点不一样，向开始介绍的那样，goalng对于client端的认证要多一个参数，生成证书的时候，要加上一个单独的认证信息：


openssl genrsa -out client.key 2048

openssl req -new -key client.key -subj "/CN=client" -out client.csr

echo extendedKeyUsage=clientAuth > extfile.cnf

openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -extfile extfile.cnf -out client.crt -days 5000 
就是多添加一个认证文件的信息，之后使用新的证书就可以实现双向认证了，这样只有那些持有被认证过的证书的客户端才能向服务端发送请求。

etcd的https的配置

docker 的https配置

k8的 apiserver的https的配置

相关参考

http://www.fenesky.com/blog/2014/07/19/how-https-works.html
http://www.ruanyifeng.com/blog/2014/02/ssl_tls.html
http://www.tuicool.com/articles/aymYbmM
