'use strict';
// Dependencies
const express = require('express');
const GitHub = require('github');
const async = require('async');
const _ = require('lodash');
const elasticsearch = require('elasticsearch');
const path = require('path');
const latex = require('../server/latex-renderer');

// Create server
const app = express();
const http = require('http').Server(app);
const io = require('socket.io')(http);

// Elasticsearch
const db = new elasticsearch.Client({
  host: 'localhost:9200',
  log: 'trace'
});

// GitHub
let github = new GitHub({
  version: '3.0.0'
});
github.authenticate({
  type: 'oauth',
  token: 'dd8fe677b60d5cbb6d0a53713821b60635caf0f5'
});
const githubPageSize = 100;

// Static: http://expressjs.com/en/starter/static-files.html
app.use(express.static('public'));
app.get('/', function(req, res) {
  res.sendfile('index.html');
});

app.use('/build', express.static(path.resolve(__dirname,'../build')));

app.get('/download/:username/cover-letter.pdf', function (req, res, next) {
  // console.log('although this matches');
  let context = req.params;
  latex.generateCoverLetter('john-doe', context)
  .then(() => {
    console.log('Done!');
    next();
  })
  .catch((error) => {
    console.log(error);
    next();
  });

  // latex.generateResume('john-doe', context)
  // .then(() => {
  //   console.log('Done!');
  // })
  // .catch((error) => {
  //   console.log(error);
  // });

})

// Socket.io API
io.on('connection', function(socket) {
  console.log('a user connected');

  socket.on('disconnect', function() {
    console.log('user disconnected');
  });

  socket.on('check-github-username', function(username, fn) {
    github.user.getFrom({
      user: username
    }, fn);
  });

  socket.on('sync-github-user', function(username, fn) {

    socket.emit('status', {
      progress: 0,
      message: 'Syncing with GitHub'
    });

    // Get User
    github.user.getFrom({
      user: username
    }, (error, user) => {
      if (error) {
        return fn(error);
      }

      // Get all Repositories
      let pages = range(1, parseInt(user.public_repos / githubPageSize) + 2, 1);
      async.map(pages, (page, cb) => {
        github.repos.getFromUser({
          user: username,
          page: page,
          per_page: githubPageSize
        }, cb);
      }, (error, repos) => {
        if (error) {
          return fn(error);
        }
        repos = _.flatten(repos);

        socket.emit('status', {
          progress: 0,
          message: 'Retrieved all repositories',
          user,
          repos
        });

        // Remove forks
        _.remove(repos, {fork:true});
        socket.emit('status', {
          progress: 0,
          message: 'Removed forked repositories'
        });

        // Get languages for Repositories
        let completedCount = 0;
        async.map(repos, (repo, cb) => {
          github.repos.getLanguages({
            user: username,
            repo: repo.name,
            per_page: githubPageSize
          }, (error, langs) => {
            delete langs.meta;
            // console.log(repo.name, langs);
            repo.languages = langs;
            completedCount++;

            socket.volatile.emit('status', {
              progress: completedCount / repos.length,
              completed: completedCount,
              total: repos.length,
              message: `Processed repository '${repo.name}'`,
              user,
              repos
            });

            return cb(error, repo);
          });
        }, (error, repos) => {

          socket.emit('status', {
            progress: 1,
            message: 'Finished analyzing repositories',
            user,
            repos
          });

          return fn(error, {
            user,
            repos
          });
        });

      });

    });
  });


  socket.on('job-search', (query='', fn) => {
    db.search({
      index: 'combined_jobs',
      size: 100,
      body: {
        query: {
          match: {
            "_all": query
          }
        }
      }
    }, (error, results) => {
      if (error) {
        return fn(error);
      }
      return fn(null, results.hits);
    });

  });

});

http.listen(3000, function() {
  console.log('listening on *:3000');
});


// ===== Helpers =======

// http://davidarvelo.com/blog/array-number-range-sequences-in-javascript-es6/
// create a generator function returning an
// iterator to a specified range of numbers
function* range(begin, end, interval = 1) {
  for (let i = begin; i < end; i += interval) {
    yield i;
  }
}