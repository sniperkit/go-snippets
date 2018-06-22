---
layout: page
title: "JavaScript - notes"
tagline: "About Electron, NW.js"
category : notes
tags : [draft, wide]
published: true
---

* [Electron - electronjs.org](https://electronjs.org/)
* [NW.js - nwjs.io](https://nwjs.io/)
    - [gitter.im](https://gitter.im) has a [Gitter desktop client](https://gitlab.com/gitlab-org/gitter/desktop/)
        - [Github (old repo.)](https://github.com/gitterHQ/desktop])
        - [Arch (AUR)](https://aur.archlinux.org/packages/gitter/)

* [Practical Node.js (online book)](https://github.com/azat-co/practicalnode)
* [What the f\*ck JavaScript? (denysdovhan/wtfjs)](https://github.com/denysdovhan/wtfjs)
* <https://github.com/ragmha/practical-js>
* <https://github.com/npm/node-semver> “The semver parser for node (the one npm uses)”

## Language notes

* `Array.prototype.push.apply(dest, [1,2,3])...` : For appending elements to an
existing array.  See also `concat()`.  Or use the new spread op. `dest.push(...src)`

* `Array.from(iterator, (x) => x+1))`

* `new Map( Array.from(iter, (e) => [ e.foo, e.bar ]) )`. It's not "easy" though
  to obtain an object from a map, have to to sthg like :

    ```javascript
    const obj = {}
    map.forEach((v,k) => { obj[k] = v })
    ```

* <https://hackernoon.com/import-export-default-require-commandjs-javascript-nodejs-es6-vs-cheatsheet-different-tutorial-example-5a321738b50f>

## Frameworks

* [2018: Reddit thread “Best Node.js web server frameworks in 2018?”](https://www.reddit.com/r/node/comments/8ajg6z/best_nodejs_web_server_frameworks_in_2018/)
  Hapi, Koa, Express, Sails, Feathers, ...?


## PWAs

* [Web App Manifest Generator @ app-manifest.firebaseapp.com/](https://app-manifest.firebaseapp.com/)
* <https://vue-starter.herokuapp.com/>
* <https://vuejsdevelopers.com/2018/04/23/vue-boilerplate-template-scaffold/>
* <https://medium.com/js-dojo/advanced-server-side-rendering-with-laravel-vue-multi-page-app-486b706e654>


### PWAs : Pointers

* [Web App Install Banners @ Google's Web Fundamentals](https://developers.google.com/web/fundamentals/app-install-banners/) __must-read__
* [Vue.js PWA template @ Github](https://github.com/vuejs-templates/pwa) [doc.](http://vuejs-templates.github.io/webpack/)
* [2017-01: We built a PWA from scratch - This is what we learned @ 14islands.com](https://14islands.com/blog/2017/01/19/progressive-web-app-from-scratch/)


### Promises

* 2012: <https://blog.domenic.me/youre-missing-the-point-of-promises/>

##### async-await lambda IIFE


```javascript
(async () => {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();
  await page.goto('https://example.com');
  await page.screenshot({path: 'example.png'});
  await browser.close();
})();
```

#### Promise.all([])

...

## Client-side

* `document.querySelectorAll("a.nav[rel=next][href^='http://example.com/abc/']")`

## Node.js

* <https://nodejs.org/api/modules.html>
* <https://github.com/babel/example-node-server> (must-read)

## ExpressJS

* [expressjs.com](https://expressjs.com)
* <http://forbeslindesay.github.io/express-route-tester/>

### Project setup

[From](https://expressjs.com/en/starter/installing.html) :

```bash
$ mkcd express-01
$ npm init
$ npm install express --save
```

#### Most basic hello world example

[From](https://expressjs.com/en/starter/hello-world.html) :

```bash
$ cat >> index.js <<EOF
const express = require('express')
const app = express()

app.get('/', (req, res) => res.send('Hello World!'))

app.listen(3000, () => console.log('Example app listening on port 3000!'))
EOF
```

#### Basic setup with React express view.

<https://github.com/reactjs/express-react-views>

```bash
$ express --view twig --css sass
```

todo...

#### Basic setup with Vue.js

<https://github.com/express-vue/express-vue>

#### React

* <https://github.com/Bogdan-Lyashenko/Under-the-hood-ReactJS>
* <https://github.com/react-boilerplate/react-boilerplate>

#### Feathers.js

* <https://github.com/feathersjs/awesome-feathersjs> resources.

#### Vue.js

* <https://madewithvuejs.com/>
* <https://madewithvuejs.com/express-vue>
    - <https://madewithvuejs.com/go/express-vue>
* <https://medium.com/@krunallathiya/vuejs-nodejs-tutorial-example-from-scratch-5a8ba947fa22>
* <https://github.com/vuejs/vue-devtools> [Chrome ext.](https://chrome.google.com/webstore/detail/vuejs-devtools/nhdogjmejiglipccpnnnanhbledajbpd)
* <https://github.com/vuejs/vuex/tree/dev/examples>

* [vue-cli](https://github.com/vuejs/vue-cli)
* [webpack template](https://vuejs-templates.github.io/webpack/)
* [Demistifying Vue's Webpack Config](https://alligator.io/vuejs/demistifying-vue-webpack/) (__must read__).

```bash
$ npm install -g vue-cli
```

```bash
$ vue init webpack vuejs-01
```

```
? Project name vuejs-01
? Project description fabic.net Vue.js experiments (01)
? Author Fabien Cadet <fabi@winterfell.fabic.net>
? Vue build standalone
? Install vue-router? Yes
? Use ESLint to lint your code? Yes
? Pick an ESLint preset Standard
? Set up unit tests Yes
```

* __TODO:__ <https://github.com/vuejs/vue-cli#customizing-webpack-config>

[Lifecycle diagram
![](https://vuejs.org/images/lifecycle.png "COPY: /assets/vuejs-instance-lifecycle.png")
](https://vuejs.org/v2/guide/instance.html#Lifecycle-Diagram)


### Compilers & static types.

* <https://babeljs.io/>

#### Static (compile-time) type checking

* <https://flow.org> (facebook's)
* <https://www.typescriptlang.org/> (microsoft)

### Webpack

* [2017: Setting Up Webpack for Bootstrap 4 and Font Awesome](https://medium.com/@estherfalayi/setting-up-webpack-for-bootstrap-4-and-font-awesome-eb276e04aaeb)
* <https://webpack.js.org/guides/shimming/>

#### Setup

* <https://github.com/postcss/autoprefixer>

### Other (todo)

* webpack, browserify
* <https://github.com/chimurai/http-proxy-middleware>

## Node.js packages

* <https://www.npmjs.com/package/bin-wrapper>

### Other

* <https://github.com/devongovett/node-wkhtmltopdf>

## Deployment

* [Process managers for Express apps @ Expressjs.com](https://expressjs.com/en/advanced/pm.html)

[install-nodejs-latest.sh](https://github.com/fabic/bash-it/blob/master/fabic/bin/setup/install-nodejs-latest.sh)

todo: setup paths

Install GCC/G++ so that NPM can build stuff :

```bash
sudo apt-get install g++
```

Using [StrongLoop's Process Manager](http://strong-pm.io/) :

```bash
npm install -g strongloop
```

* todo: [trust proxy ?](https://expressjs.com/en/guide/behind-proxies.html)

## Articles

* <https://www.html5rocks.com/en/tutorials/developertools/sourcemaps/>

## MongoDB

* [ArchLinux: MongoDB](https://wiki.archlinux.org/index.php/MongoDB)
* [MongoDB documentation](https://docs.mongodb.com/)
    - [Security & Access Control @ MongoDB manual](https://docs.mongodb.com/manual/security/)

    ```bash
    $ sudo pacman -S mongodb{,-tools}
    $ sudo systemctl start mongodb.service
    $ sudo systemctl enable mongodb.service
    ```

    And the [MongoDB shell](https://docs.mongodb.com/manual/mongo/) is
    immediately available, warning us that no access control is set up :

    ```bash
    $ mongo

    > use Scrapper
    switched to db Scrapper

    > db.Pages.insertOne({url:'http://fabic.net', label: 'fabicNet'})
    { "acknowledged" : true,
      "insertedId" : ObjectId("5ab9cf8a261bbd2cce0c71ea") }
    ```

### Setup auth/nz

» [Enable authentication (manual)](https://docs.mongodb.com/manual/tutorial/enable-authentication/)

```bash
> use admin
switched to db admin

> db.createUser({ user: 'root', pwd: 'XXXXXX', roles: [ { role: "userAdminAnyDatabase", db: "admin" } ] })

Successfully added user: {
  "user" : "root",
  "roles" : [
    {
      "role" : "userAdminAnyDatabase",
        "db" : "admin"
    }
  ]
}
```

But that is not enough, we'll need to _enable authentication_ :

```bash
$ mongo -u root -p admin

MongoDB shell version v3.6.3
connecting to: mongodb://127.0.0.1:27017
MongoDB server version: 3.6.3
2018-03-27T09:14:12.528+0400 E QUERY    [thread1] Error: Authentication failed. :
DB.prototype._authOrThrow@src/mongo/shell/db.js:1608:20
@(auth):6:1
@(auth):1:2
exception: login failed
```

```bash
> db.auth('root', 'mathinfo')

Error: Authentication failed.
```

[MongoDB configuration options](https://docs.mongodb.com/manual/reference/configuration-options/)

```bash
$ sudo journalctl -b -f
$ tail -f /var/log/mongodb/mongod.log

$ sudo systemctl stop mongodb.service

$ sudo vim /etc/mongodb.conf

## EDIT:  ##
cat >> /etc/mongodb.conf <<EOF
# /etc/mongodb.conf
# Created on $(date), by $(whoami)

systemLog:
   quiet: false
   #destination: file
   #path: "/var/log/mongodb/mongod.log"
   #logAppend: true
   destination: syslog

storage:
   dbPath: /var/lib/mongodb
   journal:
      enabled: true

net:
   bindIp: 127.0.0.1
   port: 27017

setParameter:
   enableLocalhostAuthBypass: false

security:
   authorization: enabled
EOF

$ sudo systemctl start mongodb.service
```

So, strangely, `mongo -u root -p` won't authenticate users _unless_ the
`--authenticationDatabase admin` argument is passed from the command line :

```bash
$ mongo -u root -p --authenticationDatabase admin
```

Or using `db.auth()` :

```
$ mongo

MongoDB shell version v3.6.3
connecting to: mongodb://127.0.0.1:27017
MongoDB server version: 3.6.3

> use admin
switched to db admin

> db.auth('root', 'mathinfo')
1
```

#### Create first normal user

[Security > Auth/N > Users > __Add Users__](https://docs.mongodb.com/manual/tutorial/create-users/)

```
> use fabi
switched to db fabi

> db.createUser({ user: 'fabi', pwd: 'haiku', roles: [ {role: 'dbOwner', db: 'fabi'} ] })
Successfully added user: {
  "user" : "fabi",
  "roles" : [ { "role" : "dbOwner", "db" : "fabi" } ]
}
```

```
$ mongo
> use fabi
switched to db fabi

> db.auth('fabi', 'haiku')
1
```

### Node.js

```bash
$ npm install mongodb bson-ext
```

### Random

* About __npm__'s __npx__ ?
* `$(npm bin)` typically evaluates to `./node_modules/.bin`, ex.: `$(npm bin)/cypress open`

### Pointers

* [db.collection.insertOne()](https://docs.mongodb.com/manual/reference/method/db.collection.insertOne/#db.collection.insertOne)
* [db.collection.createIndex()](https://docs.mongodb.com/manual/reference/method/db.collection.createIndex/#db.collection.createIndex)
* [db.createCollection()](https://docs.mongodb.com/manual/reference/method/db.createCollection/#db.createCollection)

__EOF__
