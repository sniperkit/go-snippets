<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Github Rank China (Dashboard star)</title>
  <meta name="viewport" content=" initial-scale=1.0, maximum-scale=1.0, user-scalable=yes">
  <link rel="stylesheet" href="style.css">
  <style>
    .avatar {
      width: 50px;
      height: 50px;
      vertical-align: middle;
    }
  </style>
  <script src="//hm.baidu.com/hm.js?4fd697a5dbf30f5d677764dc1dd3929d"></script>
</head>
<body>
<div class="container" id="root">
  <div>
    <p class="alert notice">
      <span>Last updated on </span> <strong>{{.UpdateTime}}</strong>
      <span> by Github API v3. </span>
      <span style="float:right;">
        <span>♥ made by </span>
        <a target="_blank" href="https://github.com/Piasy">Piasy</a>
        <span> just for fun, based on </span>
        <a target="_blank" href="http://githubrank.com/">githubrank.com</a>
        <span>.</span>
      </span>
    </p>
    <table class="table">
      <thead class="table-head">
      <tr>
        <th>Rank</th>
        <th>Avatar</th>
        <th>Username</th>
        <th>Name</th>
        <th>DashboardStar</th>
        <th>Followers</th>
        <th>Repos</th>
        <th>Location</th>
        <th>Updated</th>
      </tr>
      </thead>
      <tbody class="table-body">
      {{ range $rank, $user := .Ranks }}
      <tr>
        <td>
          <a name="{{$rank}}" href="#{{$rank}}">{{$rank}}</a>
        </td>
        <td>
          <img class="avatar" src="{{$user.Avatar}}&s=140">
        </td>
        <td>
          <a href="https://github.com/{{$user.Login}}" target="_blank">{{$user.Login}}</a>
        </td>
        <td>{{$user.Name}}</td>
        <td>{{$user.DashboardStar}}</td>
        <td>{{$user.Followers}}</td>
        <td>{{$user.Repos}}</td>
        <td>{{$user.Location}}</td>
        <td>{{$user.UpdatedAt}}</td>
      </tr>
      {{ end }}
      </tbody>
    </table>
  </div>
</div>
</body>
</html>
