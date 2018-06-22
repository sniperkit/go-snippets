# Changelog - Generate changelog of gitea repository

[![Build Status](https://drone.gitea.io/api/badges/go-gitea/changelog/status.svg)](https://drone.gitea.io/go-gitea/changelog)

## Purpose

This repo currently is part of Gitea. The purpose it to generate changelog when writing release notes. If a project management like Gitea, it could use this tool, otherwise please find another. The tool generate changelog depends on your PRs on the milestone and the labels of a PR.

## Installation

```
go get github.com/go-gitea/changelog
```

## Configuration

See the [changelog.yml.example](changelog.yml.example) example file.

## Usage

```
changelog -m=1.2.0 -c=/path/to/my_config_file
```

## Contributing

Fork -> Patch -> Push -> Pull Request

## Authors

* [Maintainers](https://github.com/orgs/go-gitea/people)
* [Contributors](https://github.com/go-gitea/changelog/graphs/contributors)

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/go-gitea/changelog/blob/master/LICENSE) file for the full license text.
