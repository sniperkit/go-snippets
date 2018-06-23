var gulp = require('gulp'),
  util = require('gulp-util'),
  sequence = require('run-sequence'),
  sass = require('gulp-sass'),
  uglify = require('gulp-uglify'),
  css = require('gulp-css'),
  html = require('gulp-minify-html'),
  jshint = require('gulp-jshint'),
  concat = require('gulp-concat'),
  rimraf = require('rimraf'),
  inject = require('gulp-inject'),
  rev = require('gulp-rev'),
  server = require('gulp-webserver'),
  watch = require('gulp-watch'),
  del = require('del'),
  ngAnnotate = require('gulp-ng-annotate'),
  flatten = require('gulp-flatten'),
  templateCache = require('gulp-angular-templatecache');

var distPaths = {
    root: 'dist/'
  };

distPaths.scripts = function() {
  return distPaths.root + 'js/';
};
distPaths.styles = function() {
  return distPaths.root + 'css/';
};
distPaths.images = function() {
  return distPaths.root + 'imgs/';
};
distPaths.fonts = function() {
  return distPaths.root + 'fonts/';
};

var assetFiles = {
    scripts: ['config.js',
      'app/assets/libs/angular-1*.js',
      'app/assets/libs/angular-ui-router-*.js',
      'app/assets/libs/ui-bootstrap-tpls-*.js',
      'app/assets/libs/d3-*.js'
    ],
    styles: ['app/assets/css/*.css'],
    sass: ['app/assets/sass/bootstrap/bootstrap.scss',
      'app/assets/sass/custom/*.scss'
    ],
    images: ['app/assets/imgs/**/*'],
    fonts: ['app/assets/fonts/**/*.eot',
      'app/assets/fonts/**/*.svg',
      'app/assets/fonts/**/*.ttf',
      'app/assets/fonts/**/*.woff'
    ]
  },
  appFiles = {
    modules: ['app/components/**/*.js',
      'app/sections/**/*.js',
      'app/app.js'
    ],
    templates: ['app/**/*.html',
      '!app/sections/index.html'
    ],
    index: 'app/sections/index.html'
  };

gulp.task('build:styles', function() {
  return gulp.src(assetFiles.sass.concat(assetFiles.styles))
    .pipe(sass())
    .pipe(concat('styles.css'))
    .pipe(rev())
    .pipe(gulp.dest(distPaths.styles()));
});

gulp.task('build:fonts', function() {
  return gulp.src(assetFiles.fonts)
    .pipe(flatten())
    .pipe(gulp.dest(distPaths.fonts()));
});

gulp.task('build:libs', function() {
  return gulp.src(assetFiles.scripts)
    .pipe(concat('libs.js'))
    .pipe(rev())
    .pipe(gulp.dest(distPaths.scripts()));
});

gulp.task('build:modules', function() {
  return gulp.src(appFiles.modules)
    .pipe(concat('modules.js'))
    .pipe(ngAnnotate())
    .pipe(rev())
    .pipe(gulp.dest(distPaths.scripts()));
});

gulp.task('build:views', function() {
  return gulp.src(appFiles.templates)
    .pipe(templateCache())
    .pipe(rev())
    .pipe(gulp.dest(distPaths.scripts()));
});

gulp.task('build:images', function() {
  return gulp.src(assetFiles.images)
    .pipe(flatten())
    .pipe(gulp.dest(distPaths.images()));
});


gulp.task('inject:index', function(cb) {
  return gulp.src(appFiles.index)
    .pipe(inject(gulp.src('js/*.js', {
      read: false,
      cwd: distPaths.root
    }), {
      addRootSlash: false
    }))
    .pipe(inject(gulp.src('css/*.css', {
      read: false,
      cwd: distPaths.root
    }), {
      addRootSlash: false
    }))
    .pipe(util.noop())
    .pipe(gulp.dest(distPaths.root));
});
gulp.task('webserver', function() {
  return gulp.src(distPaths.root)
    .pipe(server({
      livereload: true,
      directoryListing: false,
      open: true
    }));
});

gulp.task('remove:libs', function(cb) {
  del(distPaths.scripts() + 'libs*.js', cb);
});

gulp.task('remove:modules', function(cb) {
  del(distPaths.scripts() + 'modules*.js', cb);
});

gulp.task('remove:views', function(cb) {
  util.log('removing views at distPaths.scripts');
  del(distPaths.scripts() + 'templates*.js', cb);
  util.log('removed views');
});

gulp.task('remove:styles', function(cb) {
  del(distPaths.styles() + 'styles*.css', cb);
});

gulp.task('clean:all', function(cb) {
  rimraf(distPaths.root, cb);
});

gulp.task('start', function() {
  gulp.watch(assetFiles.sass.concat(assetFiles.styles), function(
    event) {
    utillog(event);
    sequence('remove:styles', 'build:styles', 'inject:index');
  });
  gulp.watch(assetFiles.scripts, function(event) {
    utillog(event);
    sequence('remove:libs', 'build:libs', 'inject:index');
  });
  gulp.watch(appFiles.modules, function(event) {
    utillog(event);
    sequence('remove:modules', 'build:modules', 'inject:index');
  });
  gulp.watch(appFiles.templates, function(event) {
    utillog(event);
    sequence('remove:views', 'build:views', 'inject:index');
  });
  gulp.watch(appFiles.index, function(event) {
    utillog(event);
    sequence('inject:index');
  });
});

function utillog(event) {
    util.log(util.colors.magenta(new util.File({
      path: event.path
    }).relative), event.type);
  }

gulp.task('default', function() {
  sequence('begin');
});

gulp.task('begin', function() {
  sequence('build:all', 'start', 'webserver');
});

gulp.task('build:all', ['clean:all'], function(cb) {
  sequence(['build:styles', 'build:libs', 'build:modules', 'build:views',
    'build:fonts', 'build:images'
  ], 'inject:index', cb);
});
