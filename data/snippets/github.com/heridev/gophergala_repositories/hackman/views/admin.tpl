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
		<div class="row" >
			<div class="col-md-8">
                        <blockquote>Public Announcement</blockquote>
			<form class="form-inline" action="/announce" method="POST">
				<input type="text" name="announcement" id="announcement" class="form-control" placeholder="Announcement Text" style="min-width:570px" required>
				<button type="submit" class="btn btn-default">Announce</button>
			</form>
			<hr>

                                <blockquote style="margin-top:20px;">
                                  Hackathon History
                                </blockquote>

				<table class="table table-hover" style="margin-top:20px">
					<thead>
 					<tr>
 						<th>Name</th>
 						<th>Time</th>
 					</tr>
 					</thead>
 					<tbody>
                    {{range $key, $val := .Hackathons}}
 					<tr>
 						<td><a href="/admin/hackathon/{{$val.Name}}">{{$val.Name}}</a></td>
 						<td>{{$val.Organization}}</td>
 					</tr>
                    {{end}}
 					</tbody>
				</table>
			</div>
			<div class="col-md-4">
				<blockquote>Organize Hackathon</blockquote>
				<hr>
					<form action="/organize" method="POST">
				<div class="form-group">
				    <label for="hackathon">Hackathon</label>
				    <input type="text" class="form-control" name="hackathon" id="hackathon" placeholder="Name of Hackathon" required>
				  </div>
				  <div class="form-group">
				    <label for="hackathon-description">Description</label>
				    <input type="text" class="form-control" name="hackathon-description" id="hackathon-description" placeholder="Description for Hackathon" required>
				  </div>
				  <div class="form-group">
				    <label for="hackathon-organization">Github Organization</label>
				    <input type="text" class="form-control" name="hackathon-organization" id="hackathon-organization" placeholder="GitHub Organization for Hackathon" required>
				  </div>
				  <div class="form-group">
				    <label for="hackathon-organization">Starts At</label>
				    <input type="datetime-local" class="form-control" name="start-time" id="hackathon-organization" required>
				  </div>
				  <div class="form-group">
				    <label for="hackathon-organization">Ends At</label>
				    <input type="datetime-local" class="form-control" name="end-time" id="hackathon-organization" required>
				  </div>
				<button type="submit" class="btn btn-default">Create</button>
					</form>
			</div>
		</div>
		</div>
		</div>


		<!-- <div class="col-md-4">

                        <blockquote>Pick Time And Date</blockquote>
            
				<div class="input-group">
			      <span class="input-group-btn">
			        <button class="btn btn-info" disabled="disabled" type="button">Start&nbsp;</button>
			      </span>
			      <input id="datetimepicker" type="text" type="text" class="form-control" >
			    </div>

        
		        <script>
		            jQuery('#datetimepicker').datetimepicker();
		        </script>
		        <br>
		        <div class="input-group">
			      <span class="input-group-btn">
			        <button class="btn btn-info" disabled="disabled" type="button">End&nbsp;&nbsp;</button>
			      </span>
			      <input id="datetimepicker1" type="text" type="text" class="form-control" >
			    </div>

        
		        <script>
		            jQuery('#datetimepicker1').datetimepicker();
		        </script>
		        <button class="btn btn-default" type="button" style="margin-top:10px">Announce</button>
		        <hr>
			</div> -->
			
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
