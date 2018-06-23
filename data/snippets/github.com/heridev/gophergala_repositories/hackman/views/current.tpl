<!DOCTYPE html>
<html>

	<head>
		<title>Hackman</title>
                <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
                <meta name="viewport" content="width=device-width, initial-scale=1.0">

                <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap.min.css">
	</head>

	<body>
		
		<div class="navbar navbar-default navbar-static-top">
			<div class="container">

				<div class="navbar-header">

					<a href="/" class="navbar-brand">
						<b style="color:#333; font-size:40px;">
							Hackman
						</b>
					</a>
				
					<button class="navbar-toggle" data-toggle="collapse" data-target=".navHeaderCollapse">
						<span class="icon-bar"></span>
						<span class="icon-bar"></span>
						<span class="icon-bar"></span>
					</button>
				</div>

                                <div class="collapse navbar-collapse navHeaderCollapse" >

					<ul class="nav navbar-nav navbar-right">

						<li class="dropdown">
							
							<a href="#" class="dropdown-toggle" data-toggle="dropdown">
								<img style="max-width:25px;" src="{{.Avatar}}">
                                                                {{.Name}}
								<b class="caret"></b>
							</a>
							
							<ul class="dropdown-menu" role="menu" aria-labelledby="dropdownMenuDivider">
								
								<li><a href="#">My Account</a></li>
								<li><a href="#">Github</a></li>
								<li role="presentation" class="divider"></li>
								<li ><a href="/logout">Logout</a></li>
							
							</ul>
						
						</li>
						
					</ul>
				</div>
			</div>
		</div>

		<div class="container" style="background-color:#ffffff">
		<div class="jumbotron" style="background-color:#ffffff">
                        <blockquote>Active Hackathons</blockquote>
            <table class="table">
            <tbody>

            {{range $key, $val := .Hackathons}}
            <tr>
            <td>
            <div class="row">
            	<div class="col-md-10">
            		<p style="font-size:30px">{{$val.Name}}</p>
            		<p style="font-size:14px">{{$val.Description}}</p>
            	</div>
            	<div class="col-mg-2">
            		<a href="/team?hackathonId={{$val.Id}}"><button class="btn btn-default">Participate</button></a>
            	</div>
            </div>
            </td>
            </tr>
            {{end}}
             </tbody>     
            </table>

			</div>
			
		</div>



		<div class="navbar navbar-default navbar-fixed-bottom">

			<div class="container">
				<p class="navbar-text pull-left">Made With &lt3 and Golang</p>
				<a href="#" class="navbar-btn btn-danger btn pull-right">View Project on Github</a>
			</div>
			
		</div>

		<script type="text/javascript" src="http://code.jquery.com/jquery-2.1.3.min.js"></script>
                <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/js/bootstrap.min.js"></script>

	</body>

</html>
