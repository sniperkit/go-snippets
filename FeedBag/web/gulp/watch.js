'use strict';

var gulp = require('gulp');

var paths = gulp.paths;

gulp.task('watch', ['inject'], function () {
  gulp.watch([
    paths.src + '/*.html',
    paths.src + '/{app,components}/**/*.styl',
    paths.src + '/{app,components}/**/*.js',
    paths.src + '/{app,components}/**/*.coffee',
    'bower.json'
  ], ['inject']);
});
