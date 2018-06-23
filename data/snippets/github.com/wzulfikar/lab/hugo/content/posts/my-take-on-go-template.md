---
title: "My Take on Go Template"
date: 2018-03-04T20:23:44+08:00
draft: false
tags: ["go"]
disableToc: true
---

In programming, string interpolation is common. Using ES6 (JavaScript), string interpolation can be done like this:

<!--more-->

```js
const name = "Yondu"
let greeting = `Hi ${name}!`
console.log(greeting) // yields "Hi Yondu!"
```

In Go, above code is similar to:

```go
// #1 approach: using `+`
name := "Yondu"
greeting := "Hi " + name + "!"
fmt.Println(greeting)

// #2 approach: using `fmt.Sprintf`
name := "Yondu"
greeting := fmt.Sprintf("Hi %s!", name)
fmt.Println(greeting)
```

You may have noticed that we didn't use `${...}` in Go. It's simply because that's not how we do it in Go. However, sometimes we want to interpolate more stuffs. Something like:

```
Hello world!
I'm {{ name }}. I'm {{ age }} years old!
Nice to meet you!
```

When we see above text, we may have thought of the word *"template"*, because it has a placeholder-like formatting (the use of `{{ }}`). Go provides a built-in package, [`text/template`](https://golang.org/pkg/text/template/) and [`html/template`](https://golang.org/pkg/html/template/), to deal with this template related operations.

## Enter: Go Template

![gopher](/images/gopher-head-sm.png#featured) 

Let's start with a sample code:

```go
// import "text/template"

// create template string
tplString := `
Hello world!
I'm {{ .Name }}. I'm {{ .Age }} years old!
Nice to meet you :D
`

// initialize the template
tmpl, _ := template.New("template_name").Parse(tplString)

// create the data to pass to template
data := struct {
    Name, 
    Age string
    }{
    "Yondu",
    "110",
    }

// execute the template
var tpl bytes.Buffer
_ = tmpl.Execute(&tpl, data)

// Get string representation of our template and print it
fmt.Println(tpl.String())
```
[Run in Go Playground](https://play.golang.org/p/DlSE2I-f0er)

You may notice that in above code, we used `{{ }}` in `tplString`. In Go template, the `{{ }}` is used to execute Actions. An action can mean a data evaluation or control structures. In our code, it's a data evaluation: printing variable. You should also notice the use of dot (`.`), which is how we access a variable in Go template.

In real life, you may want to separate template to its own files. For example, this is how your folder may look like:

```
▾ <mypackage>/
    ▾ templates/
        nav.tpl
        footer.tpl
        sidebar.tpl
```

And the content of those `*.tpl` files is something like:

```go
// nav.tpl
package templates

Nav := `
This is my navigation template.
{{ .Link.Home }}
{{ .Link.About }}
{{ .Link.Login }}
`
```

Using this approach – separating template files, it will be easier to edit the template and to use it. In our `nav.tpl` example above, the template string is stored in an exported variable `Nav`. An exported variable is a variable that can be accessed from other of the package. *Exporting* variable is easy: make its first letter uppercase (same goes to functions). In this case, we can access `Nav` from other package using `templates.Nav` like so:

```go
// main.go

import "mypackage/templates"

tmpl, _ := template.New("nav").Parse(templates.Nav)
```

<p class="text-center red">✽ ✽ ✽
</p>

So far, we've only used a template string from variable. Can we just pass a file and treat it as a template? Yup! You can use `.ParseFiles`:

```go
data := ... 
tmpl, _ := template.New("my-template").ParseFiles("my-template-file.txt")

var tpl bytes.Buffer
_ = tmpl.Execute(&tpl, data)

fmt.Println(tpl.String())
```

However, you may want to note that when passing a file as argument like in above `ParseFiles()` code, you should consider passing exact path. While passing relative path still works, it will be changed according to where the binary is run. So, be careful!

## Closing

It's interesting that Go provides a built-in package to handle template. Now, you can go wild and explore more about Go template (maybe, start with the difference between `text/template` and `html/template`). Have a good journey, stay safe, and don't get lost!

***Till next. See ya!***