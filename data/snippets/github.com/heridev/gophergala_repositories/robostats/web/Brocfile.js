/* global require, module */

var pickFiles = require('broccoli-static-compiler');
var EmberApp = require('ember-cli/lib/broccoli/ember-app');

var app = new EmberApp();

// Use `app.import` to add additional libraries to the generated
// output files.
//
// If you need to use different assets in different
// environments, specify an object as the first parameter. That
// object's keys should be the environment name and the values
// should be the asset to use in that environment.
//
// If the library that you are including contains AMD or ES6
// modules that you would like to import into your application
// please specify an object with the list of modules as keys
// along with the exports of each module as its value.


//c3
app.import('bower_components/c3/c3.css', {});
app.import('bower_components/d3/d3.js', {});
app.import('bower_components/c3/c3.js', {});

app.import('bower_components/bootstrap/dist/css/bootstrap.css', {});
app.import('bower_components/bootstrap-material-design/dist/css/ripples.css', {});

app.import('bower_components/bootstrap/dist/js/bootstrap.js', {});
app.import('bower_components/bootstrap-material-design/dist/js/ripples.js', {});
app.import('bower_components/bootstrap-material-design/dist/js/material.js', {});
app.import('bower_components/Chart.js/Chart.js');

app.import('bower_components/jquery-validation/dist/jquery.validate.js', {});
app.import('bower_components/jquery-validation/src/localization/messages_es.js', {});
app.import('bower_components/moment/moment.js', {});

// copy assets from packages to main app tree
var materialFont = pickFiles('bower_components/bootstrap-material-design/fonts', {
  srcDir: '/',
  files: ['Material-Design-Icons.eot', 'Material-Design-Icons.svg', 'Material-Design-Icons.ttf', 'Material-Design-Icons.woff'],
  destDir: '/assets/fonts'
});

module.exports = app.toTree(materialFont);
