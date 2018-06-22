'use strict';

const { processRequest } = require('./parse');
const readline = require('readline');

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
  terminal: false
});

rl.on('line', l => console.log(processRequest(l)));
