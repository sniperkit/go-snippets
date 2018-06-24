#!/usr/bin/env node

var http = require('http');
var os = require("os");

http.createServer(function (req, res) {
  res.writeHead(200, {'Content-Type': 'text/plain'});
  res.end('Hello World from Node.js\n');
}).listen(8080);

console.log('Server running at http://' + os.hostname() + ':8080');
