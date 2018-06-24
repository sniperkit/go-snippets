# webapp
<!-- Repo Badges for: Github Project, Slack, License-->

[![GitHub](https://img.shields.io/badge/project-Data_Together-487b57.svg?style=flat-square)](http://github.com/datatogether)
[![Slack](https://img.shields.io/badge/slack-Archivers-b44e88.svg?style=flat-square)](https://archivers-slack.herokuapp.com/)
[![License](https://img.shields.io/github/license/datatogether/webapp.svg?style=flat-square)](./LICENSE)
[![Codecov](https://img.shields.io/codecov/c/github/datatogether/webapp.svg?style=flat-square)](https://codecov.io/gh/datatogether/webapp)
[![CI](https://img.shields.io/circleci/project/github/datatogether/webapp.svg?style=flat-square)](https://circleci.com/gh/datatogether/webapp)

**Webapp** is the main public-facing web application to allow users to contribute metadata about content that has been crawled & archived from [Data Together](https://datatogether.org). A beta version of the code runs on [archivers.co](https://archivers.co) and production code on [datatogether.org](https://datatogether.org). This repository has a React front end that includes the user interface of the web application. New features of the Data Together platform are coordinated and developed in one or more of the other repositories listed in our [Roadmap](https://github.com/datatogether/roadmap). 

## License & Copyright

Copyright (C) 2017 Data Together

This program is free software: you can redistribute it and/or modify it under
the terms of the GNU Affero General Public License as published by the Free Software
Foundation, version 3.0.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE.

See the [`LICENSE`](./LICENSE) file for details.

## Getting Involved

We would love involvement from more people! If you notice any errors or would like to submit changes, please see our [Contributing Guidelines](./.github/CONTRIBUTING.md).

We use GitHub issues for [tracking bugs and feature requests](https://github.com/datatogether/webapp/issues) and Pull Requests (PRs) for [submitting changes](https://github.com/datatogether/webapp/pulls).

## Installation

To start developing, first make sure you have [Docker](https://store.docker.com/search?type=edition&offering=community) and [docker Compose](https://docs.docker.com/compose/install/). Then, run:

```shell
$ git clone git@github.com:datatogether/webapp.git
$ cd webapp
$ docker-compose up
```

This will spin up a whole bunch of servers (most of the repositories in this organization are dependencies) and make lots of noise. Once you see `webpack built` you can go to `http://localhost:3000` to see the app running locally.

## Development

Once Docker has spun up, webpack will watch `.js`, `.scss`, amd `.css` files and attempt to automatically recompile & refresh your browser window on changes. It's a good idea to have the browser's console open to make sure nothing's exploding. 

We use the following established toolkits for this project:

* [React.js](https://facebook.github.io/react/) components for all UI. If you are familiar with React, then all of the `src/js/components` files should look familiar
* [Redux Framework](https://redux.js.org/) for state control
* [React Router](https://reacttraining.com/react-router/) for navigation
* [Webpack 2](https://webpack.js.org) is our module bundler, and we also use [yarn](http://yarnpkg.com) for dependency management

### Project Organization

* `src/js/actions/` - everything the app can "do"
* `src/js/components/` - react.js components
* `src/js/containers/` - "smart" react.js components that read data from the redux state object
* `src/js/middleware/` - middleware intercepts fired-off actions, performing tasks like API calls & logging
* `src/js/propTypes/` - react.js propType definitions for reuse around the app
* `src/js/reducers/` - redux reducers that map actions to
* `src/js/selectors/` - functions that pull computed logic out of state
* `src/js/store/` - redux store configuration
* `src/js/validators/` - validation logic that checks models, returning error message validation objects
