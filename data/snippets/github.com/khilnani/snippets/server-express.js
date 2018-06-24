#!/usr/bin/env node

var os = require("os");
var express = require("express");

var app = express();

app.use(function(req, res, next){
  var ip = req.headers['x-forwarded-for'] || req.connection.remoteAddress;
  console.log('%s %s %s', req.method, req.url, ip);
  next();
});

app.get('/', function(req, res){
    res.send('Hello World! from Node.js Express');
});

app.listen(8080);
console.log('Server running at http://' + os.hostname() + ':8080');
