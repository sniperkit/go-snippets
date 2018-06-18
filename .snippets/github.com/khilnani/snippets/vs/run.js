#!/usr/bin/env node

var exec = require('child_process').exec;

function run (expr) {
  var cmd = '';

  if (expr.indexOf(".scala") > -1)
    expr = expr.split('.')[0];
  if (expr.indexOf('/') > -1)
    expr = expr.split('/').pop();

  cmd += "scala -classpath classes " + expr;
 
  console.log( "Running: " + cmd);
  exec( cmd, 
    function (error, stdout, stderr) {
      if(error)
        console.log('ERROR: ' + error)
      else 
        if( stdout)
          console.log(stdout);
    });
}

if( process.argv.length < 3 ) 
  console.log( "USAGE: " + process.argv[1].split('/').pop() + " CLASSNAME");
else
  run( process.argv[2] );
