package view

import (
	"fmt"
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/zhangjunfang/softwareScaffolding/codeGenerator/vender/lib"
)

func Write() {
	Tpl()
	for _, v := range lib.GetAllTables() {
		C(v.TableName)
		U(v.TableName)
		R(v.TableName)
		I(v.TableName)
		F(v.TableName)
	}
}

func C(tablename string) {
	u := lib.HeadToUpper(tablename)
	l := lib.OutPerfix(tablename)
	s := fmt.Sprintf(`
	<div id="content" data-dropdown="false" data-el="%s">

		<!-- %s-Create -->
		<h1>{{.Tpl.Title}}
			<div class="btn-group pull-right">
				<a href="{{urlfor "%sController.Get"}}" class="btn btn-primary">列表</a>
			</div>    
		</h1>

		<!-- %s-Form -->
		{{template "%s/_form.html" .}}

	</div>
	`, l, u, u, u, l)

	f, err := os.OpenFile(lib.Src()+lib.HeadToUpper(lib.GetDBName())+"/views/"+l+"/create.html", os.O_RDWR, 0766)
	if err != nil {
		logs.Error(err.Error())
		os.Exit(-1)
	}
	_, err = f.WriteString(s)
	if err != nil {
		logs.Error(err.Error())
		os.Exit(-1)
	}
	f.Close()
}

func U(tablename string) {
	u := lib.HeadToUpper(tablename)
	l := lib.OutPerfix(tablename)
	s := fmt.Sprintf(`
	<div id="content" data-dropdown="false" data-el="%s">

		<!-- %s-Update -->
		<h1>{{.Tpl.Title}}
			<div class="btn-group pull-right">
				<a href="{{urlfor "%sController.Get"}}" class="btn btn-primary">列表</a>
			</div>
		</h1>

		<!-- %s-Form -->
		{{template "%s/_form.html" .}}

	</div>
	`, l, u, u, u, l)

	f, err := os.OpenFile(lib.Src()+lib.HeadToUpper(lib.GetDBName())+"/views/"+l+"/update.html", os.O_RDWR, 0766)
	if err != nil {
		logs.Error(err.Error())
		os.Exit(-1)
	}
	_, err = f.WriteString(s)
	if err != nil {
		logs.Error(err.Error())
		os.Exit(-1)
	}
	f.Close()
}

func R(tablename string) {
	u := lib.HeadToUpper(tablename)
	l := lib.OutPerfix(tablename)

	var id string
	for _, v := range lib.GetAllFields(tablename) {
		if v.Key == "PRI" {
			id = v.Name
			break
		} else {
			logs.Error(tablename + ":The primary key lost")
			os.Exit(-1)
		}
	}

	var b string
	for _, v := range lib.GetAllFields(tablename) {
		if v.Key == "PRI" {
			b += fmt.Sprintf(`
<tr>
	<th>{{.Tpl.Attributes.ID}}</th>
	<td>{{.Tpl.Model.ID}}</td>
</tr>    
			`) + "\n"
		} else {
			b += fmt.Sprintf(`
<tr>
	<th>{{.Tpl.Attributes.%s}}</th>
	<td>{{.Tpl.Model.%s}}</td>
</tr>    
			`, lib.HeadToUpper(v.Name), lib.HeadToUpper(v.Name)) + "\n"
		}
	}

	s := fmt.Sprintf(`
	<div id="content" data-dropdown="false" data-el="%s">

		<!-- %s-View -->
		<h1>{{.Tpl.Title}}
			<div class="btn-group pull-right">
				<a href="{{urlfor "%sController.Get"}}" class="btn btn-primary">列表</a>
				<a href="{{urlfor "%sController.Del"}}?%s={{.Tpl.Model.ID}}" class="btn btn-danger" onclick="return confirm('Makesure Delete ID:{{.Tpl.Model.ID}}')">删除</a>
				<a href="{{urlfor "%sController.Update"}}?%s={{.Tpl.Model.ID}}" class="btn btn-info">修改</a>
			</div>
		</h1>

		<table class="table table-bordered table-hover" style="background-color: #fff">
			<tbody>
			%s   
			</tbody>
		</table>    

	</div>
	`, l, u, u, u, id, u, id, b)

	f, err := os.OpenFile(lib.Src()+lib.HeadToUpper(lib.GetDBName())+"/views/"+l+"/view.html", os.O_RDWR, 0766)
	if err != nil {
		logs.Error(err.Error())
		os.Exit(-1)
	}
	_, err = f.WriteString(s)
	if err != nil {
		logs.Error(err.Error())
		os.Exit(-1)
	}
	f.Close()
}

func I(tablename string) {
	u := lib.HeadToUpper(tablename)
	l := lib.OutPerfix(tablename)

	var id string
	for _, v := range lib.GetAllFields(tablename) {
		if v.Key == "PRI" {
			id = v.Name
			break
		} else {
			logs.Error(tablename + ":The primary key lost")
			os.Exit(-1)
		}
	}

	var b1 string
	var b2 string
	for _, v := range lib.GetAllFields(tablename) {
		if v.Key == "PRI" {
			b1 += fmt.Sprintf(`<th>{{.Tpl.Attributes.ID}}</th>`)
			b2 += fmt.Sprintf(`<td>{{$v.ID}}</td>`)
		} else {
			b1 += fmt.Sprintf(`<th>{{.Tpl.Attributes.%s}}</th>`, lib.HeadToUpper(v.Name))
			b2 += fmt.Sprintf(`<td>{{$v.%s}}</td>`, lib.HeadToUpper(v.Name))
		}
	}

	s := fmt.Sprintf(`
	<div id="content" data-dropdown="false" data-el="%s">

		<!-- %s -->
		<h1>{{.Tpl.Title}}<a href={{urlfor "%sController.Create"}} class="btn btn-primary pull-right">创建</a></h1>

		<table class="table table-bordered table-hover" style="background-color: #fff">
			<tbody>
				<tr><th>#</th>%s<th>Operate</th></tr>
				{{range $k,$v := .Tpl.Model}}
				<tr>
					<td>{{$k}}</td>%s
					<td>
						<a href="{{urlfor "%sController.View"}}?%s={{$v.ID}}"><i class="glyphicon glyphicon-eye-open"></i></a>
						<a href="{{urlfor "%sController.Update"}}?%s={{$v.ID}}"><i class="glyphicon glyphicon-pencil"></i></a>
						<a href="{{urlfor "%sController.Del"}}?%s={{$v.ID}}" onclick="return confirm('Makesure Delete ID:{{$v.ID}}')"><i class="glyphicon glyphicon-trash"></i></a>
					</td>
				</tr>
				{{end}}
			</tbody>
		</table>

	</div>
	`, l, u, u, b1, b2, u, id, u, id, u, id)

	f, err := os.OpenFile(lib.Src()+lib.HeadToUpper(lib.GetDBName())+"/views/"+l+"/index.html", os.O_RDWR, 0766)
	if err != nil {
		logs.Error(err.Error())
		os.Exit(-1)
	}
	_, err = f.WriteString(s)
	if err != nil {
		logs.Error(err.Error())
		os.Exit(-1)
	}
	f.Close()
}

func F(tablename string) {
	l := lib.OutPerfix(tablename)

	var b string
	for _, v := range lib.GetAllFields(tablename) {
		if v.Key == "PRI" {
			b += fmt.Sprintf(`
<div class="form-group">
	<label for="%s">{{.Tpl.Attributes.ID}}</label>
	<input type="text" id="%s" class="form-control" name="%s" value="{{.Tpl.Model.ID}}">
</div>   
			`, v.Name, v.Name, v.Name) + "\n"
		} else {
			b += fmt.Sprintf(`
<div class="form-group">
	<label for="%s">{{.Tpl.Attributes.%s}}</label>
	<input type="text" id="%s" class="form-control" name="%s" value="{{.Tpl.Model.%s}}">
</div>
			`, v.Name, lib.HeadToUpper(v.Name), v.Name, v.Name, lib.HeadToUpper(v.Name)) + "\n"
		}
	}

	s := fmt.Sprintf(`
	<!-- {{.Title}}-Form -->
	<form action="{{.Tpl.Action}}" class="form" method="post" enctype="multipart/form-data">
		%s
		<div class="form-group pull-right">
			<input type="reset" value="重置" class="btn btn-danger">
			<input type="submit" value="发布" class="btn btn-success">
		</div>
	</form>
	`, b)

	f, err := os.OpenFile(lib.Src()+lib.HeadToUpper(lib.GetDBName())+"/views/"+l+"/_form.html", os.O_RDWR, 0766)
	if err != nil {
		logs.Error(err.Error())
		os.Exit(-1)
	}
	_, err = f.WriteString(s)
	if err != nil {
		logs.Error(err.Error())
		os.Exit(-1)
	}
	f.Close()
}

func Tpl() {
	var b string
	for _, v := range lib.GetAllTables() {
		b += fmt.Sprintf(`<li id="%s"><a href="{{urlfor "%sController.Get"}}">%s</a></li>`, lib.OutPerfix(v.TableName), lib.HeadToUpper(v.TableName), lib.HeadToUpper(v.TableName))
	}

	s := fmt.Sprintf(`
	<html>
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>{{.Tpl.Title}}</title>
		<link rel="stylesheet" href="https://cdn.bootcss.com/bootstrap/3.3.6/css/bootstrap.min.css">
		<style>
		/*footer*/
		.footer{
			margin-top: 30px;
			text-align: center;
		}
		.footer-copy{
			font-weight: 600;
		}
		.footer-link a{
			color: #d44413;
		}
		</style>
	</head>
	<body>
	<nav class="navbar navbar-default">
		<div class="container-fluid">

			<!--Navbar-collapse and navbar-brand-->
			<div class="navbar-header">
				<button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar-collapse" aria-expanded="false">
					<span class="sr-only">Toggle navigation</span>
					<span class="icon-bar"></span>
					<span class="icon-bar"></span>
					<span class="icon-bar"></span>
				</button>
				<a class="navbar-brand" href="#">Beego</a>
			</div>

			<!--Navbar-links-->
			<div class="collapse navbar-collapse" id="navbar-collapse">
				<ul class="navbar-nav nav">
					<li id="index"><a href="{{urlfor "MainController.Get"}}">Index</a></li>
					%s
				</ul>

				<!--Navbar-widgets-->
				<form action="#" class="navbar-form navbar-right">
					<div class="form-group">
						<input type="text" name="" id="" class="form-control pull-left" placeholder="Search">
						<button class="btn btn-info pull-right">Search</button>
					</div>
				</form>

			</div>

		</div>
	</nav>

	<!--Content and Main-->
	<div class="container">

		{{if .Tpl.Success}}
		<!-- Session Message Success -->
		<div class="alert alert-success" role="alert">
			<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
			<strong>Success! </strong>{{.Tpl.Success}}
		</div>
		{{end}}
		{{if .Tpl.Danger}}
		<!-- Session Message Danger -->
		<div class="alert alert-danger" role="alert">
			<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
			<strong>Danger! </strong>{{.Tpl.Danger}}
		</div>
		{{end}}

		<!-- Content Output -->
		{{.LayoutContent}}
		
		</div>

		<!--footer-->
		<div class="footer">
			<div class="footer-copy">© BeegoCURD 2017. All rights reserved.</div>
			<div class="footer-link">
				<a href="https://my.oschina.net/u/3090506" target="_blank">OSC</a>
				<a href="https://gitee.com/Lione" target="_blank">Gitee</a>
				<a href="mailto:harikisei@live.com">Live</a>
			</div>
		</div>

		</body>
		<script src="https://cdn.bootcss.com/jquery/3.2.1/jquery.js"></script>
		<script src="https://cdn.bootcss.com/bootstrap/3.3.6/js/bootstrap.min.js"></script>
		<script>
		// 导航栏Active
		var dropdown = $('#content').attr('data-dropdown');
		var el = $('#content').attr('data-el');
		$('#'+el).addClass('active');
		if (dropdown == "true"){
			var el_next = $('#content').attr('data-el-next');
			$('#'+el_next).addClass('active');
		}
		</script>
	</html>
	`, b)
	f, err := os.OpenFile(lib.Src()+lib.HeadToUpper(lib.GetDBName())+"/views/public/tpl.html", os.O_RDWR, 0766)
	if err != nil {
		logs.Error(err.Error())
		os.Exit(-1)
	}
	_, err = f.WriteString(s)
	if err != nil {
		logs.Error(err.Error())
		os.Exit(-1)
	}
	f.Close()
}
