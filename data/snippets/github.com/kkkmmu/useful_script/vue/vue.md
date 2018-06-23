1. Vue.js 的核心是一个允许采用简洁的模板语法来声明式地将数据渲染进 DOM 的系统.
2. 组件系统是 Vue 的另一个重要概念，因为它是一种抽象，允许我们使用小型、独立和通常可复用的组件构建大型应用。仔细想想，几乎任意类型的应用界面都可以抽象为一个组件树.
3. 在 Vue 里，一个组件本质上是一个拥有预定义选项的一个 Vue 实例。在 Vue 中注册组件很简单.
4. 子单元通过 prop 接口与父单元进行了良好的解耦。
5. 在一个大型应用中，有必要将整个应用程序划分为组件，以使开发更易管理。
6. 在一个大型应用中，有必要将整个应用程序划分为组件，以使开发更易管理。在后续教程中我们将详述组件，不过这里有一个 (假想的) 例子，以展示使用了组件的应用模板是什么样的：
        ```js 
         <div id="app">
            <app-nav></app-nav>
            <app-view>
                <app-sidebar></app-sidebar>
                <app-content></app-content>
            </app-view>
         </div>
        ```
7. 每个 Vue 应用都是通过用 Vue 函数创建一个新的 Vue 实例开始的. vm (ViewModel 的缩写) 这个变量名表示 Vue 实例.
8. 一个 Vue 应用由一个通过 new Vue 创建的根 Vue 实例，以及可选的嵌套的、可复用的组件树组成。 所有的 Vue 组件都是 Vue 实例，并且接受相同的选项对象 (一些根实例特有的选项除外)。
9. 当一个 Vue 实例被创建时，它向 Vue 的响应式系统中加入了其 data 对象中能找到的所有的属性。当这些属性的值发生改变时，视图将会产生“响应”，即匹配更新为新的值。 当这些数据改变时，视图会进行重渲染。值得注意的是只有当实例被创建时 data 中存在的属性才是响应式的。
10. 除了数据属性，Vue 实例还暴露了一些有用的实例属性与方法。它们都有前缀 $，以便与用户定义的属性区分开来。
11. 每个 Vue 实例在被创建时都要经过一系列的初始化过程——例如，需要设置数据监听、编译模板、将实例挂载到 DOM 并在数据变化时更新 DOM 等。同时在这个过程中也会运行一些叫做生命周期钩子的函数，这给了用户在不同阶段添加自己的代码的机会。
    1. 比如 created 钩子可以用来在一个实例被创建之后执行代码.
    2. 也有一些其它的钩子，在实例生命周期的不同阶段被调用，如 mounted、updated 和 destroyed。生命周期钩子的 this 上下文指向调用它的 Vue 实例.
12. Vue.js 使用了基于 HTML 的模板语法，允许开发者声明式地将 DOM 绑定至底层 Vue 实例的数据。所有 Vue.js 的模板都是合法的 HTML ，所以能被遵循规范的浏览器和 HTML 解析器解析。
13. 如果你熟悉虚拟 DOM 并且偏爱 JavaScript 的原始力量，你也可以不用模板，直接写渲染 (render) 函数，使用可选的 JSX 语法。
14. 你的站点上动态渲染的任意 HTML 可能会非常危险，因为它很容易导致 XSS 攻击。请只对可信内容使用 HTML 插值，绝不要对用户提供的内容使用插值。
15. 对于所有的数据绑定，Vue.js 都提供了完全的 JavaScript 表达式支持。
16. 指令 (Directives) 是带有 v- 前缀的特殊属性。指令属性的值预期是单个 JavaScript 表达式 (v-for 是例外情况，稍后我们再讨论)。指令的职责是，当表达式的值改变时，将其产生的连带影响，响应式地作用于 DOM。
17. 一些指令能够接收一个“参数”，在指令名称之后以冒号表示。
18. 修饰符 (Modifiers) 是以半角句号 . 指明的特殊后缀，用于指出一个指令应该以特殊方式绑定。例如，.prevent 修饰符告诉 v-on 指令对于触发的事件调用 event.preventDefault().
19. v- 前缀作为一种视觉提示，用来识别模板中 Vue 特定的特性。当你在使用 Vue.js 为现有标签添加动态行为 (dynamic behavior) 时，v- 前缀很有帮助，然而，对于一些频繁用到的指令来说，就会感到使用繁琐。同时，在构建由 Vue.js 管理所有模板的单页面应用程序 (SPA - single page application) 时，v- 前缀也变得没那么重要了。因此，Vue.js 为 v-bind 和 v-on 这两个最常用的指令，提供了特定简写：
        ```javascript 
           v-bind 缩写
            <!-- 完整语法 -->
                <a v-bind:href="url">...</a>
     
            <!-- 缩写 -->
                <a :href="url">...</a>
           v-on 缩写
  
            <!-- 完整语法 -->
                <a v-on:click="doSomething">...</a>

            <!-- 缩写 -->
                <a @click="doSomething">...</a>```
        ```
20. 模板内的表达式非常便利，但是设计它们的初衷是用于简单运算的。在模板中放入太多的逻辑会让模板过重且难以维护。所以，对于任何复杂逻辑，你都应当使用计算属性。
21. 你可以像绑定普通属性一样在模板中绑定计算属性。
22. 我们可以将同一函数定义为一个方法而不是一个计算属性。两种方式的最终结果确实是完全相同的。然而，不同的是计算属性是基于它们的依赖进行缓存的。计算属性只有在它的相关依赖发生改变时才会重新求值。这就意味着只要 message 还没有发生改变，多次访问 reversedMessage 计算属性会立即返回之前的计算结果，而不必再次执行函数。相比之下，每当触发重新渲染时，调用方法将总会再次执行函数。
23. Vue 提供了一种更通用的方式来观察和响应 Vue 实例上的数据变动：侦听属性。当你有一些数据需要随着其它数据变动而变动时，你很容易滥用 watch——特别是如果你之前使用过 AngularJS。然而，通常更好的做法是使用计算属性而不是命令式的 watch 回调。
24. 计算属性默认只有 getter ，不过在需要时你也可以提供一个 setter. 
25. 虽然计算属性在大多数情况下更合适，但有时也需要一个自定义的侦听器。这就是为什么 Vue 通过 watch 选项提供了一个更通用的方法，来响应数据的变化。当需要在数据变化时执行异步或开销较大的操作时，这个方式是最有用的。
26. 操作元素的 class 列表和内联样式是数据绑定的一个常见需求。因为它们都是属性，所以我们可以用 v-bind 处理它们：只需要通过表达式计算出字符串结果即可。不过，字符串拼接麻烦且易错。因此，在将 v-bind 用于 class 和 style 时，Vue.js 做了专门的增强。表达式结果的类型除了字符串之外，还可以是对象或数组。
27. 我们可以传给 v-bind:class 一个对象，以动态地切换 class.
28. 你可以在对象中传入更多属性来动态切换多个 class。此外，v-bind:class 指令也可以与普通的 class 属性共存.
29. 我们也可以在这里绑定一个返回对象的计算属性。这是一个常用且强大的模式.
30. 我们可以把一个数组传给 v-bind:class，以应用一个 class 列表.
31. 当在一个自定义组件上使用 class 属性时，这些类将被添加到该组件的根元素上面。这个元素上已经存在的类不会被覆盖。
32. v-bind:style 的对象语法十分直观——看着非常像 CSS，但其实是一个 JavaScript 对象。CSS 属性名可以用驼峰式 (camelCase) 或短横线分隔 (kebab-case，记得用单引号括起来) 来命名.
33. v-bind:style 的数组语法可以将多个样式对象应用到同一个元素上.
34. 因为 v-if 是一个指令，所以必须将它添加到一个元素上。但是如果想切换多个元素呢？此时可以把一个 <template> 元素当做不可见的包裹元素，并在上面使用 v-if。最终的渲染结果将不包含 <template> 元素。
35. 你可以使用 v-else 指令来表示 v-if 的“else 块”, v-else 元素必须紧跟在带 v-if 或者 v-else-if 的元素的后面，否则它将不会被识别。
36. 类似于 v-else，v-else-if 也必须紧跟在带 v-if 或者 v-else-if 的元素之后。
37. 另一个用于根据条件展示元素的选项是 v-show 指令。用法大致一样：
    <h1 v-show="ok">Hello!</h1>
    不同的是带有 v-show 的元素始终会被渲染并保留在 DOM 中。v-show 只是简单地切换元素的 CSS 属性 display。
38. v-if 是“真正”的条件渲染，因为它会确保在切换过程中条件块内的事件监听器和子组件适当地被销毁和重建。
    v-if 也是惰性的：如果在初始渲染时条件为假，则什么也不做——直到条件第一次变为真时，才会开始渲染条件块。
    相比之下，v-show 就简单得多——不管初始条件是什么，元素总是会被渲染，并且只是简单地基于 CSS 进行切换。
    一般来说，v-if 有更高的切换开销，而 v-show 有更高的初始渲染开销。因此，如果需要非常频繁地切换，则使用 v-show 较好；如果在运行时条件很少改变，则使用 v-if 较好。
39. 我们用 v-for 指令根据一组数组的选项列表进行渲染。v-for 指令需要使用 item in items 形式的特殊语法，items 是源数据数组并且 item 是数组元素迭代的别名。
40. 在 v-for 块中，我们拥有对父作用域属性的完全访问权限。v-for 还支持一个可选的第二个参数为当前项的索引。
41. 你也可以用 v-for 通过一个对象的属性来迭代。你也可以提供第二个的参数为键名, 第三个参数为索引.
42. 为了给 Vue 一个提示，以便它能跟踪每个节点的身份，从而重用和重新排序现有元素，你需要为每项提供一个唯一 key 属性。理想的 key 值是每项都有的且唯一的 id。这个特殊的属性相当于 Vue 1.x 的 track-by ，但它的工作方式类似于一个属性，所以你需要用 v-bind 来绑定动态值 (在这里使用简写).
43. Vue 包含一组观察数组的变异方法，所以它们也将会触发视图更新。这些方法如下：
    push()
    pop()
    shift()
    unshift()
    splice()
    sort()
    reverse()
    变异方法 (mutation method)，顾名思义，会改变被这些方法调用的原始数组。相比之下，也有非变异 (non-mutating method) 方法，例如：filter(), concat() 和 slice() 。这些不会改变原始数组，但总是返回一个新数组。当使用非变异方法时，可以用新数组替换旧数组.
44. 对于已经创建的实例，Vue 不能动态添加根级别的响应式属性。但是，可以使用 Vue.set(object, key, value) 方法向嵌套对象添加响应式属性。
45. 类似于 v-if，你也可以利用带有 v-for 的 <template> 渲染多个元素。
46. 在自定义组件里，你可以像任何普通元素一样用 v-for 。
47. 然而，任何数据都不会被自动传递到组件里，因为组件有自己独立的作用域。为了把迭代数据传递到组件里，我们要用 props. 不自动将 item 注入到组件里的原因是，这会使得组件与 v-for 的运作紧密耦合。明确组件数据的来源能够使组件在其他场合重复使用。
48. 可以用 v-on 指令监听 DOM 事件，并在触发时运行一些 JavaScript 代码。 然而许多事件处理逻辑会更为复杂，所以直接把 JavaScript 代码写在 v-on 指令中是不可行的。因此 v-on 还可以接收一个需要调用的方法名称。除了直接绑定到一个方法，也可以在内联 JavaScript 语句中调用方法. 有时也需要在内联语句处理器中访问原始的 DOM 事件。可以用特殊变量 $event 把它传入方法。
49. 在事件处理程序中调用 event.preventDefault() 或 event.stopPropagation() 是非常常见的需求。尽管我们可以在方法中轻松实现这点，但更好的方式是：方法只有纯粹的数据逻辑，而不是去处理 DOM 事件细节。
    为了解决这个问题，Vue.js 为 v-on 提供了事件修饰符。之前提过，修饰符是由点开头的指令后缀来表示的。
        .stop
        .prevent
        .capture
        .self
        .once
50. 在监听键盘事件时，我们经常需要检查常见的键值。Vue 允许为 v-on 在监听键盘事件时添加按键修饰符.
    记住所有的 keyCode 比较困难，所以 Vue 为最常用的按键提供了别名:
        .enter
        .tab
        .delete (捕获“删除”和“退格”键)
        .esc
        .space
        .up
        .down
        .left
        .right
    可以通过全局 config.keyCodes 对象自定义按键修饰符别名：
51. 你可以用 v-model 指令在表单 <input> 及 <textarea> 元素上创建双向数据绑定。它会根据控件类型自动选取正确的方法来更新元素。尽管有些神奇，但 v-model 本质上不过是语法糖。它负责监听用户的输入事件以更新数据，并对一些极端场景进行一些特殊处理。
52. 组件 (Component) 是 Vue.js 最强大的功能之一。组件可以扩展 HTML 元素，封装可重用的代码。在较高层面上，组件是自定义元素，Vue.js 的编译器为它添加特殊功能。在有些情况下，组件也可以表现为用 is 特性进行了扩展的原生 HTML 元素。
    所有的 Vue 组件同时也都是 Vue 的实例，所以可接受相同的选项对象 (除了一些根级特有的选项) 并提供相同的生命周期钩子。
53. 要注册一个全局组件，可以使用 Vue.component(tagName, options)。组件在注册之后，便可以作为自定义元素 <my-component></my-component> 在一个实例的模板中使用。注意确保在初始化根实例之前注册组件.
54. 你不必把每个组件都注册到全局。你可以通过某个 Vue 实例/组件的实例选项 components 注册仅在其作用域中可用的组件.
55. 构造 Vue 实例时传入的各种选项大多数都可以在组件里使用。只有一个例外：data 必须是函数。
56. 组件设计初衷就是要配合使用的，最常见的就是形成父子组件的关系：组件 A 在它的模板中使用了组件 B。它们之间必然需要相互通信：父组件可能要给子组件下发数据，子组件则可能要将它内部发生的事情告知父组件。然而，通过一个良好定义的接口来尽可能将父子组件解耦也是很重要的。这保证了每个组件的代码可以在相对隔离的环境中书写和理解，从而提高了其可维护性和复用性。
    在 Vue 中，父子组件的关系可以总结为 prop 向下传递，事件向上传递。父组件通过 prop 给子组件下发数据，子组件通过事件给父组件发送消息。
57. 组件实例的作用域是孤立的。这意味着不能 (也不应该) 在子组件的模板内直接引用父组件的数据。父组件的数据需要通过 prop 才能下发到子组件中。
    子组件要显式地用 props 选项声明它预期的数据：
    Vue.component('child', {
          // 声明 props
          props: ['message'],
            // 就像 data 一样，prop 也可以在模板中使用
            // 同样也可以在 vm 实例中通过 this.message 来使用
            template: '<span>{{ message }}</span>'
   })
    然后我们可以这样向它传入一个普通字符串：
        <child message="hello!"></child>
58. 与绑定到任何普通的 HTML 特性相类似，我们可以用 v-bind 来动态地将 prop 绑定到父组件的数据。每当父组件的数据变化时，该变化也会传导给子组件.
59. 如果你想把一个对象的所有属性作为 prop 进行传递，可以使用不带任何参数的 v-bind (即用 v-bind 而不是 v-bind:prop-name)。
60. Prop 是单向绑定的：当父组件的属性变化时，将传导给子组件，但是反过来不会。这是为了防止子组件无意间修改了父组件的状态，来避免应用的数据流变得难以理解。
    另外，每次父组件更新时，子组件的所有 prop 都会更新为最新值。这意味着你不应该在子组件内部改变 prop。如果你这么做了，Vue 会在控制台给出警告。
    在两种情况下，我们很容易忍不住想去修改 prop 中数据：
        Prop 作为初始值传入后，子组件想把它当作局部数据来用；
        Prop 作为原始数据传入，由子组件处理成其它数据输出。
    对这两种情况，正确的应对方式是：
        定义一个局部变量，并用 prop 的值初始化它：
            props: ['initialCounter'],
            data: function () {
                return { counter: this.initialCounter }
            }
        定义一个计算属性，处理 prop 的值并返回：
            props: ['size'],
            computed: {
                normalizedSize: function () {
                    return this.size.trim().toLowerCase()
                }
            }
61. 我们可以为组件的 prop 指定验证规则。如果传入的数据不符合要求，Vue 会发出警告。这对于开发给他人使用的组件非常有用。
62. 所谓非 prop 特性，就是指它可以直接传入组件，而不需要定义相应的 prop。
    尽管为组件定义明确的 prop 是推荐的传参方式，组件的作者却并不总能预见到组件被使用的场景。所以，组件可以接收任意传入的特性，这些特性都会被添加到组件的根元素上。
63. 我们知道，父组件使用 prop 传递数据给子组件。但子组件怎么跟父组件通信呢？这个时候 Vue 的自定义事件系统就派得上用场了.
64. 每个 Vue 实例都实现了事件接口，即：
        使用 $on(eventName) 监听事件
        使用 $emit(eventName) 触发事件
    另外，父组件可以在使用子组件的地方直接用 v-on 来监听子组件触发的事件。
65. 在设计组合使用的组件时，内容分发 API 是非常有用的机制。
66. 通过使用保留的 <component> 元素，并对其 is 特性进行动态绑定，你可以在同一个挂载点动态切换多个组件.
67. 在编写组件时，最好考虑好以后是否要进行复用。一次性组件间有紧密的耦合没关系，但是可复用组件应当定义一个清晰的公开接口，同时也不要对其使用的外层数据作出任何假设。
    Vue 组件的 API 来自三部分——prop、事件和插槽：
        Prop 允许外部环境传递数据给组件；
        事件允许从组件内触发外部环境的副作用；
        插槽允许外部环境将额外的内容组合在组件中。
68. 在使用 vue-cli 创建的项目中，组件的创建非常方便，只需要新建一个 .vue 文件，然后在 <template> 中写好 HTML 代码，一个简单的组件就完成了. 然后在其他文件的 js 里面引入并注册，就能直接使用这个组件了
69. 一个完整的组件，除了 <template> 以外，还有 <script> 和 <style>. 需要注意的是，<script> 中的 data 必须是函数.
70. Vue 自身保留的 <component> 元素，可以将组件动态绑定到 is 特性上，从而很方便的实现动态组件切换.
71. 用 Vue.js + vue-router 创建单页应用，是非常简单的。使用 Vue.js ，我们已经可以通过组合组件来组成应用程序，当你要把 vue-router 添加进来，我们需要做的是，将组件(components)映射到路由(routes)，然后告诉 vue-router 在哪里渲染它们。
72. 该文档通篇都常使用 router 实例。留意一下 this.$router 和 router 使用起来完全一样。我们使用 this.$router 的原因是我们并不想在每个独立需要封装路由的组件中都导入路由.
73. 在vue-router中, 我们看到它定义了两个标签<router-link> 和<router-view>来对应点击和显示部分。<router-link> 就是定义页面中点击的部分，<router-view> 定义显示部分，就是点击后，区配的内容显示在什么地方。所以 <router-link> 还有一个非常重要的属性 to，定义点击之后，要到哪里去， 如：<router-link  to="/home">Home</router-link>
74. 一条路由由两部分组成： path和component.  path 指路径，component 指的是组件。如：{path:’/home’, component: home}
    我们这里有两条路由，组成一个routes: 

    const routes = [
          { path: '/home', component: Home },
          { path: '/about', component: About }
     ]
75. 创建router 对路由进行管理，它是由构造函数 new vueRouter() 创建，接受routes 参数。
    const router = new VueRouter({
              routes // routes: routes 的简写
              })
76. 配置完成后，把router 实例注入到 vue 根实例中,就可以使用路由了
    const app = new Vue({
          router
          }).$mount('#app')
77. 当用户点击 router-link 标签时，会去寻找它的 to 属性， 它的 to 属性和 js 中配置的路径{ path: '/home', component: Home}  path 一一对应，从而找到了匹配的组件， 最后把组件渲染到 <router-view> 标签所在的地方。
78. 打开浏览器控制台，首先看到 router-link 标签渲染成了 a 标签，to 属性变成了a 标签的 href 属性，这时就明白了点击跳转的意思。router-view 标签渲染成了我们定义的组件，其实它就是一个占位符，它在什么地方，匹配路径的组件就在什么地方，所以 router-link 和router-view 标签一一对应，成对出现。
79. 一个『路径参数』使用冒号 : 标记。当匹配到一个路由时，参数值会被设置到 this.$route.params，可以在每个组件内使用。
80. 你可以在一个路由中设置多段『路径参数』，对应的值都会设置到 $route.params 中。
81. 当使用路由参数时，例如从 /user/foo 导航到 /user/bar，原来的组件实例会被复用。因为两个路由都渲染同个组件，比起销毁再创建，复用则显得更加高效。不过，这也意味着组件的生命周期钩子不会再被调用。复用组件时，想对路由参数的变化作出响应的话，你可以简单地 watch（监测变化） $route 对象.
82. children 配置就是像 routes 配置一样的路由配置数组，所以呢，你可以嵌套多层路由。
83. 除了使用 <router-link> 创建 a 标签来定义导航链接，我们还可以借助 router 的实例方法，通过编写代码来实现。
    router.push(location, onComplete?, onAbort?)
    注意：在 Vue 实例内部，你可以通过 $router 访问路由实例。因此你可以调用 this.$router.push。
84. 想要导航到不同的 URL，则使用 router.push 方法。这个方法会向 history 栈添加一个新的记录，所以，当用户点击浏览器后退按钮时，则回到之前的 URL。
85. 有时候想同时（同级）展示多个视图，而不是嵌套展示，例如创建一个布局，有 sidebar（侧导航） 和 main（主内容） 两个视图，这个时候命名视图就派上用场了。你可以在界面中拥有多个单独命名的视图，而不是只有一个单独的出口。如果 router-view 没有设置名字，那么默认为 default。
    <router-view class="view one"></router-view>
    <router-view class="view two" name="a"></router-view>
    <router-view class="view three" name="b"></router-view>
    一个视图使用一个组件渲染，因此对于同个路由，多个视图就需要多个组件。确保正确使用 components 配置（带上 s）
86. For Vue, there is a parameter called e1. It takes the id of the DOM element. In the above example, we have the id #vue_det. It is the id of the div element, which is present in .html.
     <div id="vue_det"></div>
    Now, whatever we are going to do will affect the div element and nothing outside it.
87. To assign any attribute to HMTL tag, we need to use v-bind directive.
88. It's really easy to create components with Vue.js. The only three things you have to do are as follows:
    1. Create a component, and give it a template, data, methods, and whatever you need to give to it.
    2. Register it in the Vue app under the components object.
    3. Use it within the application's template.
89.  Each component can be stored in its own file with its own HTML, JavaScript, and CSS code. These are special files with the .vue extension. Inside each file, there's a <script> section for the JavaScript code, a <style> section for the CSS code, and a<template> section for the HTML code. 
90. The three most important parts of a Vuex store are state, getters, and mutations:
    State: This is an initial state of the application, basically the data of the application
    Getters: These are exactly what you think, functions that return data from the store
    Mutations: These are functions that can mutate data on the store

