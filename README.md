# Content

[![GitHub](https://img.shields.io/badge/project-Data_Together-487b57.svg?style=flat-square)](http://github.com/datatogether)
[![Slack](https://img.shields.io/badge/slack-Archivers-b44e88.svg?style=flat-square)](https://archivers-slack.herokuapp.com/)
[![License](https://img.shields.io/github/license/datatogether/content.svg)](./LICENSE)
[![Codecov](https://img.shields.io/codecov/c/github/datatogether/content.svg?style=flat-square)](https://codecov.io/gh/datatogether/content)
[![CircleCI](https://img.shields.io/circleci/project/github/datatogether/content.svg?style=flat-square)](https://circleci.com/gh/datatogether/content)

Content is a service for serving archived content stored on amazon S3.

## License & Copyright

Copyright (C) 2017 Data Together  
This program is free software: you can redistribute it and/or modify it under
the terms of the GNU Affero General Public License as published by the Free Software
Foundation, version 3.0.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE.

## Getting Involved

We would love involvement from more people! If you notice any errors or would like to submit changes, please see our [Contributing Guidelines](./.github/CONTRIBUTING.md).

We use GitHub issues for [tracking bugs and feature requests](https://github.com/datatogether/content/issues) and Pull Requests (PRs) for [submitting changes](https://github.com/datatogether/content/pulls)

## Usage

We're working on a `docker-compose` file for this guy. In the meantime, to run this service, you'll have to connect to an amazon S3 bucket and postgres server by setting the following environment variables:
```shell
  AWS_REGION
  AWS_S3_BUCKET_NAME
  AWS_S3_BUCKET_PATH
  POSTGRES_DB_URL
```

once those are properly set:
```shell
  # cd to project directory, then:
  go install

  content
```

## Development

Coming soon!
