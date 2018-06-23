//var http = require('http');
//var fs = require('fs');
//var index = fs.readFileSync('app/index.html');
//
//http.createServer(function (req, res) {
//    res.setHeader('Access-Control-Allow-Origin', '*');
//    res.setHeader('Access-Control-Request-Method', '*');
//	res.setHeader('Access-Control-Allow-Methods', 'OPTIONS, GET');
//	res.setHeader('Access-Control-Allow-Headers', '*');
//	if ( req.method === 'OPTIONS' ) {
//		res.writeHead(200);
//		res.end();
//		return;
//	}
//    res.writeHead(200, {'Content-Type': 'text/html', 'Access-Control-Allow-Origin': '*'});
//    // change the to 'text/plain' to 'text/html' it will work as your index page
//    res.end(index);
//}).listen(3001);

var express = require('express');
var cors = require('cors');
var app = express();

app.use(express.static(__dirname));

var allowCrossDomain = function(req, res, next) {
  res.header('Access-Control-Allow-Origin', '*');
  res.header('Access-Control-Allow-Methods', 'GET,PUT,POST,DELETE');
  res.header('Access-Control-Allow-Headers', 'X-Requested-With, Accept, Origin, Referer, User-Agent, Content-Type, Authorization');
 
  // intercept OPTIONS method
  if (req.method === 'OPTIONS') {
    res.send(200);
  }
  else {
    next();
  }
};

  // you might have "app" instead of "server"
//app.configure(function() {
//    app.use(allowCrossDomain);   // make sure this is is called before the router
//    app.use(app.router);      // not entirely necessary--will be automatically called with the first .get()
//});

//app.use(allowCrossDomain);
app.use(function(req, res, next) {
  res.header("Access-Control-Allow-Origin", "*");
  res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
  next();
});
app.listen(process.env.PORT || 3001);
