使用Golang的Context管理上下文
golang 1.7版本中context库被很多标准库的模块所使用，比如net/http和os的一些模块中，利用这些原生模块，我们就不需要自己再写上下文的管理器了，直接调用函数接口即可实现，利用context我们可以实现一些比如请求的声明周期内的变量管理，执行一些操作的超时等等。

保存上下文对象

这里我们通过一个简单的例子来看一下如何使用context的特性来实现上下文的对象保存，这里我们写了一个简单的http server，具有登录和退出，状态检查路由（检查用户是否登录）

func main(){
mux:=http.NewServeMux()
		mux.HandleFunc("/",StatusHandler)
		mux.HandleFunc("/login",LoginHandler)
		mux.HandleFunc("/logout",LogoutHandler)
		contextedMux:=AddContextSupport(mux)
		log.Fatal(http.ListenAndServe(":8080",contextedMux))
}
其中的AddContextSupport是一个中间件，用来绑定一个context到原来的handler中，所有的请求都必须先经过该中间件后才能进入各自的路由处理中。具体的实现代码如下：

func AddContextSupport(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method, "-", r.RequestURI)
			cookie, _ := r.Cookie("username")
			if cookie != nil {
ctx := context.WithValue(r.Context(), "username", cookie.Value)
// WithContext returns a shallow copy of r with its context changed
// to ctx. The provided ctx must be non-nil.
next.ServeHTTP(w, r.WithContext(ctx))
} else {
next.ServeHTTP(w, r)
}
})
}
该中间件可以打印每次请求的方法和请求的url，然后获得请求的cookie值，如果cookie为空的话则继续传递到对应的路由处理函数中，否则保存cookie的值到Context, 注意这里的Context()是request对象的方法, 将创建一个新的上下文对象（如果context为空），context.WithValue()函数将key和value保存在新的上下文对象中并返回该对象。

其余的路由处理函数代码如下, 分别用来保存cookie的登录路由LoginHandler()，还有删除cookie的退出路由处理函数LogoutHandler()。

func LoginHandler(w http.ResponseWriter,r *http.Request){
expitation := time.Now().Add(24*time.Hour)
				var username string
				if username=r.URL.Query().Get("username");username==""{
					username = "guest"
				}
cookie:=http.Cookie{Name:"username",Value:username,Expires:expitation}
	   http.SetCookie(w,&cookie)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
expiration := time.Now().AddDate(0, 0, -1)
				cookie := http.Cookie{Name: "username", Value: "alice_cooper@gmail.com", Expires: expiration}
			http.SetCookie(w, &cookie)
}
这里我们在请求/login的时候，可以传递一个参数username到函数中，比如/login?username=alice , 默认为”guest”用户. 设置的cookie有效期为1天，删除的时候我们只需要设置cookie为之前的日期即可。

另外一个关键的部分是读取上下文保存内容的 StatusHandler() 路由处理函数，该函数将调用r.Context()获得request的上下文，如果我们执行了login后，那我们在中间件中已经设置了该对象，所以请求将查看是否上下文对象中保存了一个名为username的对象，如果有的话则回应一个欢迎页面。否则告知用户没有登录。

func StatusHandler(w http.ResponseWriter,r *http.Request){

	if username:=r.Context().Value("username"); username!=nil{
		w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hi username:"+username.(string)+"\n"))
	}else{
		w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Logged in"))
	}
}
我们不仅仅可以在上下文中保存简单的类型，我们可以保存任何类型的数据，因为Value() 返回的对象是一个interface{}对象，所以我们需要转换一下才能使用。

超时处理

对于简单的保存和传递对象，使用context的确很方便，但是该库的使用不仅仅是保存变量，还可以创建一个超时和取消的行为，比如说我们web端去请求了其他的资源，但是该资源的处理比较耗时，我们无法预见什么时候能够返回，如果让用户超时的话，实在是不太好，所以我们需要创建一个超时的操作，主动判断是否超时，然后传递一个合适的行为给用户。

这里我们现在路由中增加一个长期运行的job路由

mux.HandleFunc("/longjob",jobWithCancelHandler)
具体的处理如下，我们的handler会利用WithCancel() 返回一个新的（如果没有创建）或者原来已保存的上下文，还有一个cancel对象，这个对象可以用来手动执行取消操作。另外我们的url中可以指定这个任务模拟执行的长度，比如/longjob?jobtime=10则代表模拟的任务将会执行超过10秒。 执行任务的函数longRunningCalculation（）返回一个chan该chan会在执行时间到期后写入一个Done字符串。

handler中我们就可以使用select语句监听两个非缓存的channel，阻塞直到有数据写到任何一个channel中。比如代码中我们设置了超时是5秒，而任务执行10秒的话则5秒到期后ctx.Done()会因为cancel()的调用而写入数据，这样该handler就会因为超时退出。否则的话则执行正常的job处理后获得传递的“Done”退出。

func longRunningCalculation(timeCost int)chan string{

result:=make(chan string)
		   go func (){
			   time.Sleep(time.Second*(time.Duration(timeCost)))
				   result<-"Done"
		   }()
	   return result
}
func jobWithCancelHandler(w http.ResponseWriter, r * http.Request){
	var ctx context.Context
		var cancel context.CancelFunc
		var jobtime string
		if jobtime=r.URL.Query().Get("jobtime");jobtime==""{
			jobtime = "10"
		}
	timecost,err:=strconv.Atoi(jobtime)
		if err!=nil{
			timecost=10
		}
	log.Println("Job will cost : "+jobtime+"s")
		ctx,cancel = context.WithCancel(r.Context())
		defer cancel()

		go func(){
			time.Sleep(5*time.Second)
				cancel()
		}()

	select{
		case <-ctx.Done():
			log.Println(ctx.Err())
				return
		case result:=<-longRunningCalculation(timecost):
					io.WriteString(w,result)
	}
	return
}

这就是使用context的一些基本方式，其实context还有很多函数这里没有涉及，比如WithTimeout和WithDeadline等，但是使用上都比较相似，这里就不在举例。

另外：我的个人网站jsmean.com欢迎大家访问。
