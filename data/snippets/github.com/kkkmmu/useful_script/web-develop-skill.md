#Tools
1. Chrome :F12打开Debug 模式，点击Network就可以看到页面的交互流程以及每一步交互的数据（包括异步获取的，ajax, js）. 同时通过Chrome的Debug 模式Elements选项可以查看动态生成页面的最终DOM文件。
2. Phantomjs 很有用
3. Casperjs 很有用
4. golang html/template的分割符“{{}}”可以通过 template.Template的func (t *Template) Delims(left, right string) *Template来重定义。通过这种方法可以解决golang template与一些前端框架（比如说vue）的分割符冲突问题。
5. vue.js的template 分割符“{{}}”可以通过vue.Look 或者vue.config.delimiters=XXX来重新定义。这可以用来解决分割符冲突问题。
6.  一个典型的golang web server: (js/css/image/file通过FileServer来实现). 
	```http.HandleFunc("/", MainPage)
	http.Handle("/asset/web/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":8080", nil)```

7. 注意golang的http.Handle中路由表项末尾带"/"和不带"/"其结果是截然不同的.
8. 要在给客户端的回复中带上Cookie可以 通过http.SetCookie来设置)
#9  Web APP的开发要尽量做到前后端分离，也就是说server端不应该对前端的表示做任何假设，同样前端也不应该对server端的实现做任何假设。前端通过访问服务端定义的一组良好的接口来获取和提交数据，（这也是前端和后端应该仅有的交互）
10. chrome的console中可以直接修改页面目前的参数,以及脚本中的变量.
11. vue // 模板定界符 修改
	Vue.config.delimiters = ['{{', '}}']

	// html模板定界符
	Vue.config.unsafeDelimiters = ["{{{", "}}}"]

12. 要将前端的表单数据转换成json(JSON.stringify())以后再传到后端，这样就会使后端的表单解析得到极大的简化。
13.要想将页脚固定到页面最底部需要先设置页面高度，然后再设置footer高度，参考mra.css 和footer.css
14. VUE常用第三方库：
	1. vue-router
	2. vuex
	3. element
	4. vux
	5. iview

15. VUE 单文件组件:
	在很多Vue项目中，我们使用 Vue.component 来定义全局组件，紧接着用 new Vue({ el: '#container '}) 在每个页面内指定一个容器元素。
	这种方式在很多中小规模的项目中运作的很好，在这些项目里 JavaScript 只被用来加强特定的视图。但当在更复杂的项目中，或者你的前端完全由 JavaScript 驱动的时候，下面这些缺点将变得非常明显：
	全局定义(Global definitions) 强制要求每个 component 中的命名不得重复
	字符串模板(String templates) 缺乏语法高亮，在 HTML 有多行的时候，需要用到丑陋的 \
	不支持CSS(No CSS support) 意味着当 HTML 和 JavaScript 组件化时，CSS 明显被遗漏
	没有构建步骤(No build step) 限制只能使用 HTML 和 ES5 JavaScript, 而不能使用预处理器，如 Pug (formerly Jade) 和 Babel
	文件扩展名为 .vue 的 single-file components(单文件组件) 为以上所有问题提供了解决方法，并且还可以使用 Webpack 或 Browserify 等构建工具。
16: Jquery ajax 获取response Status Code:
	$.get("example.url.com", function(data) {
		console.log(data);
		}).done(function() {
			// TO DO ON DONE
		}).fail(function(data, textStatus, xhr) {
			//This shows status code eg. 403
			console.log("error", data.status);
			//This shows status message eg. Forbidden
			console.log("STATUS: "+xhr);
		}).always(function() {
			//TO-DO after fail/done request.
			console.log("ended");
		});
17. 使用Flex进行页面布局。
18. 通过chrome的console可以查看和修改javacript中定义的变量，对象。
