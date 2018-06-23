# Stacker

Stacker is a container manager for [Docker][]

This is currently beta software. You can find me as `cnf` on the channel `#docker` on `irc.freenode.net`

pull requests, especially for documentation, are welcome. For new features, open a ticket or poke me on irc.

`docker pull frosquin/stacker:dev` gets you the dev version, which is mostly up to date with the devel branch.

`make image` will build a docker image of your own.

You can find the [full documentation here](https://github.com/cnf/stacker/wiki)

## Configuration

### toml files

Stacker uses [toml][] files for its configuration

```toml
[stacker]
config_dir = "conf.d"
watch = true

[docker]
socket = "/var/run/docker.sock"

[<module>]
option = "foo"
```

#### \[stacker\] options
  * `config_dir`:string - where to look for container config files.
  * `watch`:bool - to watch for and react to changes in files in `config_dir`

#### \[docker\] options
  * `socket`:string - the docker socket to connect to.
    - `http:/example.host.local:2375/`
    - `unix:///var/run/docker.sock`

TODO:
  * `cert`:string
  * `key`:string
  * `ca`:string

#### planned modules
  * [consul](http://consul.io) - for configuration and service advertising
  * [etcd](https://github.com/coreos/etcd) - for configuration and service advertising
  * `logger` and `syslog` - for local and remote logging of container actions.
  * `cron` - for scheduled runs of containers.

#### \[\[container\]\] options
  * `name`:string - name of the container. This is a mandatory field.
  * `hostname`:string - hostname of the container.
  * `user`:string - user to run as in the container.
  * `memory`:string - memory limits to place on the container. minimum of 512K.
  * `cpu_shares`:int - relative weight of CPU usage vs other containers.
  * `cpu_set`:string - CPU's in which to allow execution (0-3, 0,1)
  * `attach`:list - attach to `stdin`, `stdout`, `stderr`
  * `expose`:list - expose a port from the container. This does not publish it to the host.
  * `tty`:bool - allocate a pseudo-tty.
  * `env`:list - list of environment variables. `["FOO=bar", "FOO2=baz"]`
  * `cmd`:list - docker CMD list.
  * `image`:string - docker image to use.
  * `volumes`:list - list of volumes.
  * `workdir`:string - working directory inside the container.
  * `entrypoint`:list - overwrite the default entrypoint set by the image.

  * `cap_add`:list - add Linux capabilities.
  * `cap_drop`:list - drop Linux capabilities.
  * `cid_file`:string - write the container ID to the file.
  * `lxc_conf`:list - add custom lxc options.
  * `privileged`:bool - give extended privileges to this container.
  * `publish`:list - publish a container᾿s port to the host. (format: ip:hostPort:containerPort | ip::containerPort | hostPort:containerPort | containerPort)
  * `publish_all`:bool - publish all exposed ports to the host interfaces.
  * `link`:list - add link to another container (name:alias)
  * `dns`:list - set custom dns servers for the container.
  * `dns_search`:list - set costom dns search domain for the container.
  * `volumes_from`:list - mount all volumes from the given container(s).
  * `net`:string - set the Network mode for the container. (`bridge`, `none`, `container:<name|id>`, `host`)  [ref](http://docs.docker.com/reference/run/#network-settings)
  * `remove`:bool - are we allowed to remove this container? NOTE: this differs from the `docker run` option.

## TODO

Lots :P

[Docker]: http://docker.com
[toml]: https://github.com/toml-lang/toml
