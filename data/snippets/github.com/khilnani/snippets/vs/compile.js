#!/usr/bin/env node

var exec = require('child_process').exec;

function compile (expr) {
  console.log( "Preparing...");

  var cmd = 'rm -rf classes && mkdir classes && ';

  if(expr == undefined || expr == null)
    cmd += "scalac -d classes src/*.scala";
  else if (expr.indexOf(".scala") > -1)
    cmd += "scalac -d classes " + expr ;
  else
    cmd += "scalac -d classes " + expr + ".scala" 

  cmd += " -feature" 
 
  console.log( "Compiling: " + cmd);
  exec( cmd, 
    function (error, stdout, stderr) {
      if(error)
        console.log('ERROR: ' + error)
      else 
        if( stdout)
          console.log(stdout);
      console.log( "Done.");
    });
}

if( process.argv.length < 3 ) {
  console.log( "No File specified. Using *");
  compile();
} else
  compile( process.argv[2] );
