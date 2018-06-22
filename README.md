# HashiCorp Vault on Alpine

A super-simple Vault image that is more suitable for use as a client than the [official image](https://hub.docker.com/_/vault/).

No `VOLUME`s nor `EXPOSE`s -- *no side-effects*. Just an `ENTRYPOINT` to the `vault` executable.

```
$ docker run -it dweomer/vault

Usage: vault <command> [args]

Common commands:
    read        Read data and retrieves secrets
    write       Write data, configuration, and secrets
    delete      Delete secrets and configuration
    list        List data or secrets
    login       Authenticate locally
    server      Start a Vault server
    status      Print seal and HA status
    unwrap      Unwrap a wrapped secret

Other commands:
    audit          Interact with audit devices
    auth           Interact with auth methods
    lease          Interact with leases
    operator       Perform operator-specific tasks
    path-help      Retrieve API help for paths
    policy         Interact with policies
    secrets        Interact with secrets engines
    ssh            Initiate an SSH session
    token          Interact with tokens
```
