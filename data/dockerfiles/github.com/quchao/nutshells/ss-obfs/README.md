# SS-OBFS Docker Image

[![Project Nutshells](https://img.shields.io/badge/Project-_Nutshells_🌰-orange.svg)](https://github.com/quchao/nutshells/) [![Docker Build Build Status](https://img.shields.io/docker/build/nutshells/ss-obfs.svg?maxAge=3600&label=Build%20Status)](https://hub.docker.com/r/nutshells/ss-obfs/) [![Alpine Based](https://img.shields.io/badge/Alpine-3.7-0D597F.svg)](https://alpinelinux.org/) [![CentOS Based](https://img.shields.io/badge/CentOS-7.5-932279.svg)](https://www.centos.org/) [![MIT License](https://img.shields.io/github/license/quchao/nutshells.svg?label=License)](https://github.com/quchao/nutshells/blob/master/LICENSE) [![SS-libev](https://img.shields.io/badge/shadowsocks--libev-3.2.0-lightgrey.svg)](https://github.com/shadowsocks/shadowsocks-libev/) [![Simple-obfs](https://img.shields.io/badge/Simple--obfs-0.0.5-lightgrey.svg)](https://github.com/shadowsocks/simple-obfs)

This image combines [ss-libev](https://github.com/shadowsocks/shadowsocks-libev) and [simple-obfs](https://github.com/shadowsocks/simple-obfs) both with the server and client included.


## Variants:

| Tag | Description | 🐳 |
|:-- |:-- |:--:|
| `:latest` | Shadowsocks-libev & simple-obfs on `alpine:latest`. | [![Dockerfile](https://img.shields.io/badge/Dockerfile-latest-22B8EB.svg?maxAge=2592000&style=flat-square)](https://github.com/quchao/nutshells/blob/master/ss-obfs/Dockerfile/) |
| `:centos` | Shadowsocks-libev & simple-obfs on `centos:latest`. | [![Dockerfile](https://img.shields.io/badge/Dockerfile-latest-22B8EB.svg?maxAge=2592000&style=flat-square)](https://github.com/quchao/nutshells/blob/master/ss-obfs/Dockerfile-centos/) |


## Usage

### Synopsis

```
docker container run [OPTIONS] nutshells/ss-obfs [COMMAND] [ARG...]
```

- Learn more about `docker container run` and its `OPTIONS` [here](https://docs.docker.com/edge/engine/reference/commandline/container_run/);
- List all available `COMMAND`s: 
    `docker container run --rm --read-only nutshells/ss-obfs help`
- List all `ARG`s (for [configuring](#custom-configurations) the containers):
    - For *server* mode:
    `docker container run --rm --read-only nutshells/ss-obfs server --help`
    - For *client* mode:
    `docker container run --rm --read-only nutshells/ss-obfs client --help`

### Getting Started

Let's create a *server* container on the remote server:

> Replace `<password>` with a [strong](https://en.wikipedia.org/wiki/Password_strength) password, and assign an unused port to `<remote_port>`.

``` bash
docker container run \
       -d -p <remote_port>:12345/tcp -p <remote_port>:12345/udp \
       --name=ss-obfs-server --restart=unless-stopped --read-only \
       -e PASSWORD=<password> \           
       nutshells/ss-obfs
       server
```

Then run a *client* container on your local machine with the same configuration:

> Remove `127.0.0.1:` from the `-p` options if you want to share it with the local network.

``` bash
docker container run \
       -d -p 127.0.0.1:12345:12345/tcp -p 127.0.0.1:12345:12345/udp \
       --name=ss-obfs-client --restart=unless-stopped --read-only \
       -e PASSWORD=<password> \
       -e SERVER_ADDRESS=<remote_address> -e SERVER_PORT=<remote_port> \           
       nutshells/ss-obfs
       client
```

See if it shows the correct IP address of your remote server.

``` bash
curl -sSL https://ipinfo.io/ --socks5 127.0.0.1:12345
```

### Custom Configurations

As you can see from the examples of the previous sections: the container accepts original [command-line options](https://github.com/shadowsocks/shadowsocks-libev/#usage) of *ss-libev* as arguments.

Here are a few more examples:

#### Enabling verbose mode:

``` bash
docker container run [OPTIONS] nutshells/ss-obfs server -v [OTHER_ARG...]
```

#### Printing the version & command-line options:

``` bash
docker container run --rm --read-only nutshells/ss-obfs server --help
docker container run --rm --read-only nutshells/ss-obfs client --help
```

However, please be informed that **some** of the options are managed by [the entrypoint script](https://github.com/quchao/nutshells/blob/master/ss-obfs/docker-entrypoint.sh) of the container. You will encounter an error while trying to set any of them, just use the [environment variables](#environment-variables) instead; as for other exceptions, just follow the message to get rid of them.


## Reference

### Environment Variables

| Name | Default | Relevant Option | Description |
|:-- |:-- |:-- |:-- |
| `SERVER_ADDRESS` |  | `-s` | Hostname or IP address of the remote server. **Required** for `client` mode only. |
| `SERVER_PORT` |  | `-p` | Port number of the remote server. **Required** for `client` mode only. |
| `PASSWORD` |  | `-p` | Password of the remote server. **Required** if `KEY_IN_BASE64` is unset. |
| `KEY_IN_BASE64` |  | `--key` | Base64 encoded password of the remote server. **Required** if `PASSWORD` is unset. |
| `ENCRYPT_METHOD` | `xchacha20-ietf-poly1305` | `-m` | Encryption Algorithms. Modern ciphers such as `chacha20` or `aes-gcm` based are always recomended;  choose from the latter ones if the server is AES-NI enabled. |
| `TCP_RELAY` | `true` | `-u` & `-U` | Enable TCP relay. |
| `UDP_RELAY` | `true` | `-u` & `-U` | Enable UDP relay. |
| `REUSE_PORT` | `true` | `--reuse-port` | Enable port reuse. |
| `TCP_FAST_OPEN` | `true` | `--fast-open` | Enable TCP fast open. |
| `OBFS_PLUGIN` | `http` | `--plugin` | Traffic obscuration plugin. `http` or `tls`. |
| `OBFS_HOST` | `bing.com` | `--plugin-opts` | Hostname for traffic obscuration. |

Learn to set an environment variable [here](https://docs.docker.com/engine/reference/commandline/run/#set-environment-variables--e-env-env-file), or refer to the [Getting Started](#getting-started) section for some live examples.

### Build Arguments

| Argument | Description | Default |
|:-- |:-- |:--:|
| `WITH_OBFS` | Compile with the [simple-obfs](https://github.com/shadowsocks/simple-obfs) plugin. | `true` |


## Advanced Topics

### Using Docker-compose

See [the sample file](https://github.com/quchao/nutshells/blob/master/ss-obfs/docker-compose.yml).

### Better Performance

According to the [optimization tips](https://github.com/shadowsocks/shadowsocks/wiki/Optimizing-Shadowsocks), you'd better increase the open file limit for the container:

> It's NOT recommended to enable the privileged mode to this purpose.
> The default limit for open files is `1024` in containers.

``` bash
docker container run \
       --ulimit nofile=51200:51200 \
       [OTHER_OPTIONS] nutshells/ss-obfs server [OTHER_ARG...]
```

Additionally, the `--net=host` option provides the best network performance, use it if you know it exactly.

### Gaining a shell access

Get an interactive shell to a **running** container:

> Run `/bin/bash` instead in the CentOS variants.

``` bash
docker container exec -it <container_name> /bin/ash
```

### Customizing the image

#### by using the `--build-arg` options

You can customize the image more easily by using the `--build-arg` option.

> Check out all the **build-time** arguments in [this table](#build-arguments).

Let's remove the simple-obfs plugin completely from the build: 

``` bash
git clone --depth 1 https://github.com/quchao/nutshells.git
docker image build -q=false --rm=true --no-cache=true \
             -t nutshells/ss-obfs \
             -f ./ss-obfs/Dockerfile \
             --build-arg WITH_OBFS=false \
             ./ss-obfs
```

#### By modifying the dockerfile

You may want to make some modifications to the image.
Pull the source code from GitHub, customize it, then build one by yourself:

``` bash
git clone --depth 1 https://github.com/quchao/nutshells.git
docker image build -q=false --rm=true --no-cache=true \
             -t nutshells/ss-obfs \
             -f ./ss-obfs/Dockerfile \
             ./ss-obfs
```

#### By committing the changes on a container

Otherwise just pull the image from the official registry, start a container and [get a shell](#gaining-a-shell-access) to it, [commit the changes](https://docs.docker.com/engine/reference/commandline/commit/) afterwards.

``` bash
docker container commit --change "Commit msg" <container_name> nutshells/ss-obfs
```


## Contributing

[![Github Starts](https://img.shields.io/github/stars/quchao/nutshells.svg?maxAge=3600&style=social&label=Star)](https://github.com/quchao/nutshells/) [![Twitter Followers](https://img.shields.io/twitter/follow/chappell.svg?maxAge=3600&style=social&label=Follow)](https://twitter.com/chappell/)

> Follow GitHub's [*How-to*](https://opensource.guide/how-to-contribute/) guide for the basis.

Contributions are always welcome in many ways:

- Give a star to show your fondness;
- File an [issue](https://github.com/quchao/nutshells/issues) if you have a question or an idea;
- Fork this repo and submit a [PR](https://github.com/quchao/nutshells/pulls);
- Improve the documentation.


## Todo

- [ ] Add a section about modifying kernel parameters in containers.
- [ ] Add a `HealthCheck` instruction.

## Acknowledgments & Licenses

Unless specified, all codes of **Project Nutshells** are released under the [MIT License](https://github.com/quchao/nutshells/blob/master/LICENSE).

Other relevant softwares:

| Ware/Lib | License |
|:-- |:--:|
| [Docker](https://www.docker.com/) | [![License](https://img.shields.io/github/license/moby/moby.svg?maxAge=2592000&label=License)](https://github.com/moby/moby/blob/master/LICENSE/) |
| [Shadowsocks-libev](https://github.com/shadowsocks/shadowsocks-libev/) | [![License](https://img.shields.io/github/license/shadowsocks/shadowsocks-libev.svg?maxAge=2592000&label=License)](https://github.com/shadowsocks/shadowsocks-libev/blob/master/LICENSE) |
| [Simple-obfs](https://github.com/shadowsocks/simple-obfs/) | [![License](https://img.shields.io/github/license/shadowsocks/simple-obfs.svg?maxAge=2592000&label=License)](https://github.com/shadowsocks/simple-obfs/blob/master/LICENSE) |
