# API Server

[![Build Status](https://travis-ci.org/DevMine/api-server.png?branch=master)](https://travis-ci.org/DevMine/api-server)
[![GoDoc](http://godoc.org/github.com/DevMine/api-server?status.svg)](http://godoc.org/github.com/DevMine/api-server)
[![GoWalker](http://img.shields.io/badge/doc-gowalker-blue.svg?style=flat)](https://gowalker.org/github.com/DevMine/api-server)

`api-server` is a JSON RESTful API server implemented in
[Go](http://golang.org/).

It serves data collected and processed by various DevMine projects
([crawld](http://devmine.ch/doc/crawld/),
[featscomp](http://devmine.ch/doc/featscomp/),
[repotool](http://devmine.ch/doc/repotool/)...).

## API Documentation

### General Information

All data is sent and received as JSON.

Timestamps use the ISO 8601 format:

```
YYYY-MM-DDThh:mm:ssTZD (eg 2015-01-09T14:19:47+01:00)
```

Blank fields are included as `null`.

Only GET requests are answered.

#### Client Errors

Sending invalid JSON will result in a `400 Bad Request` response:

```
GET /search/{"foo":1
```

***Response***

```
{
  "message": "invalid JSON input"
}
```

`400 Bad Request` response may be returned when sending valid JSON that the
server does not know how to process. Example:

```
GET /search/{"foo":1}
```

***Response***

```
{
  "message": "non existing feature: foo"
}
```

#### Common Parameters

Parameters not specified as a segment in the path can be passed as an HTTP query
string parameter.

Requests that return multiple results limit to 30 items by default.
However, up to 100 items can be returned by specifying the `?per_page`
parameter. Example:

```
GET /users?per_page=42
```

Since not all resources are shown on a page, further items may be queried
by specifying the `?since` parameter, which corresponds to an item ID:

```
GET /users?per_page=42&since=3747
```

### Version

All requests receive the version 1 of the API. You can verify which version of
the API the server is serving by visiting the route `/`:

```
GET /
```

***Response***

```
{
  "version": 1,
  "doc_url": "http://devmine.ch/doc/api-server"
}
```

### Users

Users related resources are served under the `/users` routes.

#### Get all users

The `/users` route provides a dump of all the users, sorted by user IDs.
As the number of results is limited, you can specify from which user ID you
would like to list the users with the `?since` parameter.

```
GET /users
```

#### Get a single user

You can get a single user by querying the `/users/:username` route.

```
GET /users/Rolinh
```

***Response***

```
{
  "id": 38769,
  "username": "Rolinh",
  "name": "Robin Hahling",
  "email": "robin.hahling@gw-computing.net",
  "gh_user": {
    "id": 38769,
    "github_id": 1324157,
    "login": "Rolinh",
    "bio": null,
    "blog": "http://projects.gw-computing.net",
    "company": "HGdev",
    "email": "robin.hahling@gw-computing.net",
    "hireable": false,
    "location": "Switzerland",
    "avatar_url": "https://avatars.githubusercontent.com/u/1324157?v=3",
    "html_url": "https://github.com/Rolinh",
    "followers_count": 8,
    "following_count": 19,
    "collaborators_count": null,
    "created_at": "2012-01-12T09:37:19+01:00",
    "updated_at": "2015-01-09T18:36:56+01:00",
    "gh_organizations": [
    {
      "id": 2522,
      "github_id": 6969061,
      "login": "DevMine",
      "avatar_url": "https://avatars.githubusercontent.com/u/6969061?v=3",
      "html_url": "https://github.com/DevMine",
      "name": "DevMine",
      "company": null,
      "blog": "http://devmine.ch/",
      "location": "Around the world",
      "email": null,
      "collaborators_count": null,
      "created_at": "2014-03-16T22:07:05+01:00",
      "updated_at": "2015-01-09T21:51:06+01:00"
    }
    ]
  }
}
```

#### Get commits of a user

You can get all the commits of a user by querying the `/users/:username/commits`
route. As there is potentially a lot of commits, results are paginated.

```
GET /users/Rolinh/commits
```

***Response***

```
[
  {
    "id": 1375919,
    "repository": {
      "id": 93271,
      "name": "crawld",
      "primary_language": "Go",
      "clone_url": "https://github.com/DevMine/crawld.git",
      "clone_path": "go/devmine/crawld",
      "vcs": "git"
    },
    "message": "crawld: Make sure we finish writing logs before exiting.\n\nLog output is buffered and written periodically using flush (around\nevery 10 seconds).\n",
    "author": {
      "id": 46138,
      "username": "Rolinh",
      "name": "Robin Hahling",
      "email": "robin.hahling@gw-computing.net"
    },
    "committer": {
      "id": 46138,
      "username": "Rolinh",
      "name": "Robin Hahling",
      "email": "robin.hahling@gw-computing.net"
    },
    "author_date": "2015-01-08T00:31:22+01:00",
    "commit_date": "2015-01-08T00:31:22+01:00",
    "file_changed_count": 1,
    "insertions_count": 3,
    "deletions_count": 0
  },
...
]
```

#### Get repositories associated to a user

You can get the repositories associated to a user by querying the
`/users/:username/repositories` route.

```
GET /users/Rolinh/repositories
```

***Response***

```
[
{
  "id": 76947,
  "name": "crawld",
  "primary_language": "Go",
  "clone_url": "https://github.com/DevMine/crawld.git",
  "clone_path": "go/devmine/crawld",
  "vcs": "git",
  "gh_repository": {
    "id": 76941,
    "github_id": 28636035,
    "full_name": "DevMine/crawld",
    "description": "A data crawler and repository fetcher",
    "homepage": "http://devmine.ch/doc/crawld/",
    "fork": false,
    "default_branch": "master",
    "master_branch": null,
    "html_url": "https://github.com/DevMine/crawld",
    "forks_count": 0,
    "open_issues_count": 1,
    "stargazers_count": 0,
    "subscribers_count": 3,
    "watchers_count": 0,
    "size_in_kb": 260,
    "created_at": "2014-12-30T16:44:02+01:00",
    "updated_at": "2015-01-09T18:37:41+01:00",
    "pushed_at": "2015-01-09T16:57:28+01:00"
  }
},
...
]
```

#### Get features scores of a user

You can get the features scores of a user by querying the
`/users/:username/scores` route.

```
GET /users/Rolinh/scores
```

***Response***

```
{
  "contributions_count": 0.48484848484848486,
  "followers_count": 0.02478026651545222,
  "forks_avg": 0.27997405412506565,
  "hireable": 1,
  "stars_avg": 0.31201177610713027
}
```

### Repositories

Repositories related resources are served under the `/repositories` routes.

##### Get all repositories

The `/repositories` route provides a dump of all the repositories, sorted by
repositories IDs.
As the number of results limited, you can specify from which repository ID you
would like to list the repositories with the `?since` parameter.

```
GET /repositories
```

#### Get repositories by name

You can get repositories by name using the `/repositories/:name` route.
Note that several repositories may have the same name. Hence, a list of
repositories is returned.

```
GET /repositories/crawld
```

***Response***

```
[
{
  "id": 76947,
  "name": "crawld",
  "primary_language": "Go",
  "clone_url": "https://github.com/DevMine/crawld.git",
  "clone_path": "go/devmine/crawld",
  "vcs": "git",
  "gh_repository": {
    "id": 76941,
    "github_id": 28636035,
    "full_name": "DevMine/crawld",
    "description": "A data crawler and repository fetcher",
    "homepage": "http://devmine.ch/doc/crawld/",
    "fork": false,
    "default_branch": "master",
    "master_branch": null,
    "html_url": "https://github.com/DevMine/crawld",
    "forks_count": 0,
    "open_issues_count": 1,
    "stargazers_count": 0,
    "subscribers_count": 3,
    "watchers_count": 0,
    "size_in_kb": 260,
    "created_at": "2014-12-30T16:44:02+01:00",
    "updated_at": "2015-01-09T18:37:41+01:00",
    "pushed_at": "2015-01-09T16:57:28+01:00"
  }
},
...
]
```

### Features

Features related resources are served under the `/features` routes.

##### Get all features

The `/features` route provides a dump of all the features, sorted by
features IDs.
As the number of results limited, you can specify from which feature ID you
would like to list the features with the `?since` parameter.

```
GET /features
```

#### Get features by feature category

Features are classified into categories. You can get a dump of all features from
a category using the `/features/by_category/:name` route.

```
GET /features/by_category/other
```

***Response***

```
[
{
  "id": 13,
  "name": "followers_count",
  "category": "other",
  "default_weight": 1
},
{
  "id": 14,
  "name": "hireable",
  "category": "other",
  "default_weight": 1
},
...
]
```

#### Get users scores by feature

The `/features/:name/scores` route provides a list of users and scores for the
given feature name.

As the number of results limited, you can specify from which user ID you would
like to list the results with the `?since` parameter.

```
GET /followers_count/scores
```

***Response***

```
[
{
  "id": 234,
  "username": "austinheap",
  "score": 0.00022682166146867
},
{
  "id": 235,
  "username": "javierprovecho",
  "score": 0.00056705415367168
},
...
]
```

### Search queries

Search queries can be done under the `/search/:query` route.

`query` is a JSON formatted input object of feature name with their weights.

Example query:

```
GET /search/{"followers_count":4}
```

The results is a list of users with their ranks, sorted from higher ranked to
lower ranked user according to the query.
The search results is limited to the top 1000 ranked users.

***Response***

```
[
  {
    "id": 2290,
    "username": "defunkt",
    "name": "Chris Wanstrath",
    "email": "chris@github.com",
    "rank": 4.292859984017694
  },
  {
    "id": 2374,
    "username": "mojombo",
    "name": "Tom Preston-Werner",
    "email": null,
    "rank": 4.095068589338459
  },
  {
    "id": 22682,
    "username": "paulirish",
    "name": "Paul Irish",
    "email": "",
    "rank": 3.211172019935161
  },
...
]
```

### Stats

Querying the `/stats` route provides some statistics about the items in the
database.

```
GET /stats
```

***Response***

```
{
  "users_count": 59171,
  "repositories_count": 121898,
  "commits_count": 10331903,
  "commit_deltas_count": 94117385,
  "features_count": 6,
  "gh_users_count": 59170,
  "gh_organizations_count": 3985,
  "gh_repositories_count": 121896
}
```

## Installation

To install the API server, run this command in a terminal, assuming
[Go](http://golang.org/) is installed:

```
go get github.com/DevMine/api-server
```

Or you can download a binary for your platform from the DevMine project's
[downloads page](http://devmine.ch/downloads).

You also need to setup a [PostgreSQL](http://www.postgresql.org/) database.
And of course, you need to add some data into your database and compute the
features (see [crawld](http://devmine.ch/doc/crawld/),
[features](http://devmine.ch/doc/features/) and other DevMine projects for
this).

Some matrix computation is done and it uses the
[BLAS](http://www.netlib.org/blas/) library so you need to have it installed on
the server as well.

## Usage and configuration

Copy `devmine.conf.sample` to `devmine.conf` and edit it according to your
needs. The configuration file has two sections:

* **database**: allows you to configure access to your PostgreSQL
  database.
  - **hostname**: hostname of the machine.
  - **port**: PostgreSQL port.
  - **username**: PostgreSQL user that has access to the database.
  - **password**: password of the database user.
  - **dbname**: database name.
  - **ssl\_mode**: takes any of these 4 values: "disable",
    "require", "verify-ca", "verify-null". Refer to PostgreSQL
    [documentation](http://www.postgresql.org/docs/9.4/static/libpq-ssl.html)
    for details.
* **server**: allows you to configure the server parameters.
  - **hostname**: server hostname.
  - **port**: port on which to listen.
  - **enable\_cors**: boolean indicating whether to allow Cross Origin Resource
    Sharing (CORS) or not.

Once the configuration file has been adjusted, you are ready to run the API
server (`devmine`).
You need to specify the path to the configuration file with the help of the `-c`
option. Example:

    devmine -c devmine.conf

Some command line options are also available, mainly about logging options.
