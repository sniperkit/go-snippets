# linterhub-docker

This is a curated catalog of linters (engines) wrapped in docker containers. Each container is based on [Alpine Linux](https://alpinelinux.org) with preinstalled engine and a minimalistic set of dependencies.

## Install

The following software is required:

- [Docker](https://www.docker.com/)

It's recommended to get containers directly from [Docker Hub](https://hub.docker.com/u/linterhub/):

```bash
docker pull linterhub/your_faivorite_engine
```

You can also build it from source:

```bash
git clone https://github.com/repometric/linterhub-docker
cd engines/your_faivorite_engine
docker build -t linterhub/your_faivorite_engine
```

## Usage

You will need to prepare engine-specific command (for analysis) and mount a local folder to container:
```bash
docker run -v /local:/container -it linterhub/your_faivorite_engine engine_command
```

## Contribute

Contributions are always welcome! Please read the [contribution guidelines](.github/CONTRIBUTING.md) first.

- [Request to add new engine](.github/ENGINE_TEMPLATE.md)
- [Report an issue for existing engine](.github/ISSUE_TEMPLATE.md)
- [Propose pull request](.github/PULL_REQUEST_TEMPLATE.md)
- [How to add new engine](.github/HOWTO_ENGINE.md)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE) file for details.
