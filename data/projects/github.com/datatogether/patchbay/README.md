# Patchbay

[![GitHub](https://img.shields.io/badge/project-Data_Together-487b57.svg?style=flat-square)](http://github.com/datatogether)
[![Slack](https://img.shields.io/badge/slack-Archivers-b44e88.svg?style=flat-square)](https://archivers-slack.herokuapp.com/)
[![License](https://img.shields.io/github/license/datatogether/patchbay.svg?style=flat-square)](./LICENSE)
[![Codecov](https://img.shields.io/codecov/c/github/datatogether/patchbay.svg?style=flat-square)](https://circleci.com/gh/datatogether/patchbay) 
[![CI](https://img.shields.io/circleci/project/github/datatogether/patchbay.svg?style=flat-square)](https://codecov.io/gh/datatogether/patchbay)

Patchbay is a websockets-based API backend for our react & redux-based [webapp](https://github.com/datatogether/webapp). Patchbay carries react _actions_ onto the server, allowing duplexed communication between webapp clients & the react server. The server is able to fire native react actions that redux will handle. We use this for reporting the progress of backend services like the completion of [task_mgmt tasks](https://github.com/datatogether/task_mgmt) via a redis pub-sub connection, leave open the option for cross-user interactions using patchbay as the coordinating server.

We structure this entire websocket integration as a protocol upgrade above the [api service](https://github.com/datatogether/api). From the webapp perspective, it should be possible to have patchbay explode in a glorious ball of fire, and fall back to using api endpoints without interruption to the user.

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

We would love involvement from more people! If you notice any errors or would like 
to submit changes, please see our [Contributing Guidelines](./CONTRIBUTING.md). 

We use GitHub issues for [tracking bugs and feature requests](https://github.com/datatogether/patchbay/issues) 
and Pull Requests (PRs) for [submitting changes](https://github.com/datatogether/patchbay/pulls)

## Usage 

Because patchbay supports actions defined in [webapp](https://github.com/datatogether/webapp), the two repos are tightly coupled. The best way to see & work with patchbay is to run it as a backing service for webapp. 

## Development

Webapp's [development instructions](https://github.com/datatogether/webapp/#development) are a great way to get started!
