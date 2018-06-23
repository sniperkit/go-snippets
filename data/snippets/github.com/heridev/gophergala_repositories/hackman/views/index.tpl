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
							HackMan
						</b>
					</a>
				
					<button class="navbar-toggle" data-toggle="collapse" data-target=".navHeaderCollapse">
						<span class="icon-bar"></span>
						<span class="icon-bar"></span>
						<span class="icon-bar"></span>
					</button>
				</div>

				<!--<div class="collapse navbar-collapse navHeaderCollapse" >

					<ul class="nav navbar-nav navbar-right">
						
						<li><a href="login.htm">Login</a></li>
						
					</ul>
				</div>-->				
			</div>
		</div>
		<div class="container">
			<div class="jumbotron">
				<p>Sign in with GitHub</p>
				<hr>
				<div class="form-group">
                                <a class="btn btn-primary" href="https://github.com/login/oauth/authorize?client_id={{.ClientID}}&scope=user:email" role="button">User</a>
                                <a class="btn btn-primary" href="https://github.com/login/oauth/authorize?client_id={{.ClientID}}&scope=read:org,admin:org,admin:org_hook,user:email" role="button">Admin</a>
				<!--<button type="button" class="btn btn-primary">Sign In With GitHub</button><br>-->
			</div>
		</div>




		<div class="navbar navbar-default navbar-fixed-bottom">

			<div class="container">
				<p class="navbar-text pull-left">Made With &lt3 and Golang</p>
				<a href="https://github.com/pravj/beehub" class="navbar-btn btn-danger btn pull-right">View Project on Github</a>
			</div>
			
		</div>

                <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/js/bootstrap.min.js"></script>

	</body>

</html>
