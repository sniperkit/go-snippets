# Deprecated user [Viper](http://spf13.com/project/viper)
[![GoDoc](https://godoc.org/github.com/abhiyerra/goconfigure?status.svg)](https://godoc.org/github.com/abhiyerra/goconfigure)

## goconfigure - Go Library to get configuration from different sources

### Introduction

Golang configuration parser via the following ways:

 - [X] ENV
 - [X] etcd
 - [ ] JSON
 - [ ] YAML
 - [ ] flags
 - [ ] ZooKeeper
 - [ ] Consul

### Example

    etcdClient := etcd.NewClient([]string{"http://127.0.0.1:4001"})
    config := goconfigure.Config{
        Etcd: goconfigure.Etcd{
            Client:    etcdClient,
            Namespace: "goconfigure",
        },
    }
    config.Add("database_url",
        "Use etcd for configuration.",
        "user=ayerra dbname=goconfigure_example sslmode=disable",
        "GOCONFIGURE_DATABASE_URL",
        "database_url")
    config.SetParser(goconfigure.EnvParser)

    dbUrl, err := config.Get("database_url")

To get the value from Etcd change the parser to EtcdParser. The key
that is searched for etcd is /goconfigure/database_url. That is it
takes the from /:namespace/:key.

All values are returned as string.
