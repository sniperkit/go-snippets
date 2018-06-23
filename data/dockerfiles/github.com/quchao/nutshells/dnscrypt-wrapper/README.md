# DNSCrypt-Wrapper Docker Image

[![Project Nutshells](https://img.shields.io/badge/Project-_Nutshells_🌰-orange.svg?maxAge=2592000)](https://github.com/quchao/nutshells/) [![Docker Build Build Status](https://img.shields.io/docker/build/nutshells/dnscrypt-wrapper.svg?maxAge=3600&label=Build%20Status)](https://hub.docker.com/r/nutshells/dnscrypt-wrapper/) [![Alpine Based](https://img.shields.io/badge/Alpine-3.6-0D597F.svg?maxAge=2592000)](https://alpinelinux.org/) [![MIT License](https://img.shields.io/github/license/quchao/nutshells.svg?maxAge=2592000&label=License)](https://github.com/quchao/nutshells/blob/master/LICENSE) [![dnscrypt-wrapper](https://img.shields.io/badge/DNSCrypt--Wrapper-0.3-lightgrey.svg?maxAge=2592000)](https://github.com/cofyc/dnscrypt-wrapper/)

[DNSCrypt-Wrapper](https://github.com/cofyc/dnscrypt-wrapper/) is the server-end of [DNSCrypt](http://dnscrypt.org/) proxy, which is a protocol to improve DNS security, now with xchacha20 cipher support.
This image features certficate management & rotation.


## Variants:

| Tag | Description | 🐳 |
|:-- |:-- |:--:|
| `:latest` | DNSCrypt-Wrapper `0.3` on `alpine:latest`. | [![Dockerfile](https://img.shields.io/badge/Dockerfile-latest-22B8EB.svg?maxAge=2592000&style=flat-square)](https://github.com/quchao/nutshells/blob/master/dnscrypt-wrapper/Dockerfile/) |


## Usage

### Synopsis

```
docker container run [OPTIONS] nutshells/dnscrypt-wrapper [COMMAND] [ARG...]
```

- Learn more about `docker container run` and its `OPTIONS` [here](https://docs.docker.com/edge/engine/reference/commandline/container_run/);
- List all available `COMMAND`s: 
    `docker container run --rm --read-only nutshells/dnscrypt-wrapper help`
- List all `ARG`s:
    `docker container run --rm --read-only nutshells/dnscrypt-wrapper --help`

### Getting Started

A DNScrypt proxy server cannot go without a so-called **provider key pair**.
It's rather easy to generate a new pair by running the `init` [command](#commands) as below:

> `<keys_dir>` is the host directory where you store the key pairs.
> The provider key pair should NEVER be changed as you may inform the world of the public key, unless the secret key is compromised.

``` bash
docker container run -d -p 5353:12345/udp -p 5353:12345/tcp \
       --name=dnscrypt-server --restart=unless-stopped --read-only \
       --mount=type=bind,src=<keys_dir>,dst=/usr/local/etc/dnscrypt-wrapper \
       nutshells/dnscrypt-wrapper \
       init
```

Now, a server is initialized with [default settings](#environment-variables) and running on port `5353` as a daemon.

Dig the *public key* fingerprint out of logs:

``` bash
docker container logs dnscrypt-server | grep --color 'Provider public key: '
```

Then [see if it works](#how-to-test).

### Using an existing key pair

If you used to run a DNScrypt proxy server and have been keeping the keys securely, make sure they are put into `<keys_dir>` and renamed to `public.key` and `secret.key`, then run the `start` [command](#commands) instead.

> Incidentally, `start` is the default one which could be just omitted.

Lost the *public key* fingerprint but couldn't find it from logs? Try the `pubkey` command:

``` bash
docker container run --rm --read-only \
       --mount=type=bind,src=<keys_dir>,dst=/usr/local/etc/dnscrypt-wrapper \
       nutshells/dnscrypt-wrapper \
       pubkey
```

### How to test

Get `<provider_pub_key>` by following the instructions above,
and check [this section](#environment-variables) to determin the default value of `<provider_basename>`.

> Please install [`dnscrypt-proxy`](https://hub.docker.com/r/nutshells/dnscrypt-proxy/) and `dig` first.

``` bash
dnscrypt-proxy --local-address=127.0.0.1:53 \
               --resolver-address=127.0.0.1:5353 \
               --provider-name=2.dnscrypt-cert.<provider_basename> \
               --provider-key=<provider_pub_key>
               --loglevel=7
dig -p 53 +tcp google.com @127.0.0.1
```

### Utilities

As you can see from the examples of the previous sections: the container accepts original [command-line options](https://github.com/shadowsocks/shadowsocks-libev/#usage) of *ss-libev* as arguments.

Here are a few more examples:

#### Printing the version:

``` bash
docker container run --rm --read-only nutshells/dnscrypt-wrapper --version
```

#### Enabling verbose mode:

``` bash
docker container run [OPTIONS] nutshells/dnscrypt-wrapper [COMMAND] [ARG...] -V
```

#### Printing command-line options:

``` bash
docker container run --rm --read-only nutshells/dnscrypt-wrapper --help
```

However, please be informed that **some** of the options are managed by [the entrypoint script](https://github.com/quchao/nutshells/blob/master/dnscrypt-wrapper/docker-entrypoint.sh) of the container. You will encounter an error while trying to set any of them, just use the [environment variables](#environment-variables) instead; as for other exceptions, just follow the message to get rid of them.


## Reference

### Environment Variables

| Name | Default | Relevant Option | Description |
|:-- |:-- |:-- |:-- |
| `RESOLVER_IP` | `8.8.8.8` | `-r`, `--resolver-address` | Hostname or IP address of the upstream dns resolver. |
| `RESOLVER_PORT` | `53` | `-r`, `--resolver-address` | Port number of the upstream dns resolver. |
| `PROVIDER_BASENAME` | `example.com` | `--provider-name` | Basename of the provider, which forms the whole provide name with a prefix `2.dnscrypt-cert.`. |
| `CRYPT_KEYS_LIFESPAN` | `365` | `--cert-file-expire-days` | For how long (in days) the crypt key & certs would be valid. Refer to [this topic](#rotating-the-crypt-key-and-certs) to automate the rotation. |

For instance, if you want to use [OpenDNS](https://www.opendns.com) as the upstream DNS resolver other than [Google's Public DNS](https://developers.google.com/speed/public-dns/), the default one, just [set an environment variable](https://docs.docker.com/engine/reference/commandline/run/#set-environment-variables--e-env-env-file) like this:

```
docker container run \
       -e RESOLVER_IP=208.67.222.222 -e RESOLVER_PORT=5353 \
       [OTHER_OPTIONS] ...
```

### Data Volumes

| Path in Container | Description | Mount as Writeable |
|:-- |:-- |:--:|
| `/usr/local/etc/dnscrypt-wrapper` | Directory where keys are stored | Y |


## Advanced Topics

### Using Docker-compose

See [the sample file](https://github.com/quchao/nutshells/blob/master/dnscrypt-wrapper/docker-compose.yml).

### Better Performance

It's a good idea to speed up the high-frequency queries by adding a caching upstream resolver, we choose [dnsmasq](http://www.thekelleys.org.uk/dnsmasq/doc.html) for its lightweight and configurability. Though there're many other options, such as [unbound](https://www.unbound.net/).

Firstly, let's create a new bridge network named `dnscrypt`:

> Add `--driver overlay` if it's in swarm mode.

``` bash
docker network create dnscrypt
```

Secondly, create a *dnsmasq* container (as an *upstream* resolver) into the network with an increased cache size:

> Refer to [this page](https://hub.docker.com/r/nutshells/dnsmasq-fast-lookup/) for more about the `nutshells/dnsmasq-fast-lookup` image.

``` bash
docker container run -d --network=dnscrypt --network-alias=upstream \
       --name=upstream --restart=unless-stopped --read-only \
       nutshells/dnsmasq-fast-lookup \
       --domain-needed --bogus-priv \
       --server=8.8.8.8 --no-resolv --no-hosts \
       --cache-size=10240
```

Then start a *dnscrypt server* into the same network too:

> To add an existing container into the network, use [`docker network connect`](https://docs.docker.com/engine/userguide/networking/#user-defined-networks) please.

``` bash
docker container run -d --network=dnscrypt \
       -p 5353:12345/udp -p 5353:12345/tcp \
       --name=dnscrypt-server --restart=unless-stopped --read-only \
       --mount=type=bind,src=<keys_dir>,dst=/usr/local/etc/dnscrypt-wrapper \
       -e RESOLVER_IP=upstream -e RESOLVER_PORT=12345 \
       nutshells/dnscrypt-wrapper \
       init
```

Done!

A [compose file](#using-docker-compose), which helps you to manage both of the containers concurrently, is highly recommended for this situation.

Alternatively, the `--net=host` option provides the best network performance, use it if you know it exactly.

### Backing-up the secret key

If you forgot to mount `<keys_dir>` into the container,
and now you want to locate the secret key in the anonymous volume,
just do some inspection first:

``` bash
docker container inspect -f '{{json .Mounts }}' dnscrypt-server | grep --color '"Source":'
```

Then backup it securely.

### Rotating the crypt key and certs

Unlike the lifelong provider key pair, a **crypt key** & two certs, which are time-limited and used to encrypt and authenticate DNS queries, will be generated only if they're missing or expiring on starting. Thus the container is supposed to be restarted before certs' expiration.

> Two certs are issued right after the crypt key's generation, one of them uses xchacha20 cipher.

Let's say we're planning to rotate them about once a week.
Firstly, shrink [the cert's lifespan](#environment-variables) to `7` days:

> Actually the rotation starts when the validity remaining is under `30%`, which would be on day `5` in this case.

```
docker container run -e CRYPT_KEYS_LIFESPAN=7 ...
```

Secondly, restart the container every single day by creating a daily cronjob:

``` bash
0 4 * * * docker container restart dnscrypt-server
```

### Gaining a shell access

Get an interactive shell to a **running** container:

``` bash
docker container exec -it dnscrypt-server /bin/ash
```

### Customizing the image

#### By modifying the dockerfile

You may want to make some modifications to the image.
Pull the source code from GitHub, customize it, then build one by yourself:

``` bash
git clone --depth 1 https://github.com/quchao/nutshells.git
docker image build -q=false --rm=true --no-cache=true \
             -t nutshells/dnscrypt-wrapper \
             -f ./dnscrypt-wrapper/Dockerfile \
             ./dnscrypt-wrapper
```

#### By committing the changes on a container

Otherwise just pull the image from the official registry, start a container and [get a shell](#gaining-a-shell-access) to it, [commit the changes](https://docs.docker.com/engine/reference/commandline/commit/) afterwards.

``` bash
docker container commit --change "Commit msg" dnscrypt-server nutshells/dnscrypt-wrapper
```


## Caveats

## Declaring the Health Status

Status of this container-specified health check merely indicates whether the crypt certs are *about to expire*, you'd better **restart** the container to [rotate the keys](#rotating-the-crypt-key-and-certs) ASAP if it's shown as *unhealthy*.

To confirm the status, run this command:

``` bash
docker container inspect --format='{{json .State.Health.Status}}' dnscrypt-server
```

And to check the logs:

``` bash
docker container inspect --format='{{json .State.Health}}' dnscrypt-server | python -m json.tool
```

If you think this is annoying, just add [the `--no-healthcheck` option](https://docs.docker.com/engine/reference/run/#healthcheck) to disable it.


## Contributing

[![Github Starts](https://img.shields.io/github/stars/quchao/nutshells.svg?maxAge=3600&style=social&label=Star)](https://github.com/quchao/nutshells/) [![Twitter Followers](https://img.shields.io/twitter/follow/chappell.svg?maxAge=3600&style=social&label=Follow)](https://twitter.com/chappell/)

> Follow GitHub's [*How-to*](https://opensource.guide/how-to-contribute/) guide for the basis.

Contributions are always welcome in many ways:

- Give a star to show your fondness;
- File an [issue](https://github.com/quchao/nutshells/issues) if you have a question or an idea;
- Fork this repo and submit a [PR](https://github.com/quchao/nutshells/pulls);
- Improve the documentation.


## Todo

- [x] Serve with the old key & certs for another hour after the rotation.
- [x] Add instructions on how to speed it up by caching the upstream dns queries.
- [x] Add a `HealthCheck` instruction to indicate the expiration status of certs.
- [ ] Add a command for checking the expire status.
- [ ] Use another container to rotate the keys.


## Acknowledgments & Licenses

Unless specified, all codes of **Project Nutshells** are released under the [MIT License](https://github.com/quchao/nutshells/blob/master/LICENSE).

Other relevant softwares:

| Ware/Lib | License |
|:-- |:--:|
| [Docker](https://www.docker.com/) | [![License](https://img.shields.io/github/license/moby/moby.svg?maxAge=2592000&label=License)](https://github.com/moby/moby/blob/master/LICENSE) |
| [DNSCrypt-Proxy](https://github.com/jedisct1/dnscrypt-proxy) | [![License](https://img.shields.io/badge/License-ISC_License-blue.svg?maxAge=2592000)](https://github.com/jedisct1/dnscrypt-proxy/blob/master/COPYING) |
| [DNSCrypt-Wrapper](https://github.com/cofyc/dnscrypt-wrapper/) | [![License](https://img.shields.io/badge/License-ISC_License-blue.svg?maxAge=2592000)](https://github.com/cofyc/dnscrypt-wrapper/blob/master/COPYING) |
| [DNSCrypt-Server-Docker](https://github.com/jedisct1/dnscrypt-server-docker/) | [![License](https://img.shields.io/github/license/jedisct1/dnscrypt-server-docker.svg?maxAge=2592000&label=License)](https://github.com/jedisct1/dnscrypt-server-docker/blob/master/LICENSE) |
