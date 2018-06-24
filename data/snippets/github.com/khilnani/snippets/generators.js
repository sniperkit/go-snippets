#!/usr/bin/env node --harmony-generators

var fs = require('fs');


function thread(fn) {
  // fn is not executed until gen.next is called.
  var gen = fn();

  // callback for fs.readFile
  // uses closure to access the main generator to call next yield
  function next(err, res) {
    console.log('get yield.');

    // send the return data to the prior yield, let it continue exec and get the next yield
    var ret = gen.next(res);

    // value: result of initial exec of the fn, read(path) or size(path)
    if (ret.done) return;

    // effectively: read(path)(next)
    ret.value(next);
  }
  
  // kick off (err and res are undefined , but its ok)
  next();
}

thread(function *(){

  try {
    console.log('start');

    var a = yield read('a.txt');
    console.log('done with a');

    var b = yield size('a.txt');
    console.log(b);
    console.log(a)
  } catch (e) {
    console.log('Error: ' + e.message);
  }
});

function read(path) {
  var fn = function(done){
    fs.readFile(path, 'utf8', done);
  }
  return fn; 
}


function size(path) {
  var fn = function(done){
    fs.stat(path, function(err, stat) {
      done(err, stat.size + 'b'); 
    }); 
  }
  return fn; 
}
