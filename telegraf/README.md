Telegraf is collector for InfluxDB and other time series database.

[hub.docker.com/gianarb/telegraf](https://hub.docker.com/r/gianarb/telegraf/)

```bash
docker pull gianarb/telegraf:0.10.4.1
docker run -it gianarb/telegraf:0.10.4.1
```

This image use an environment variable `$TELEGRAF_CONF` (default /etc/telegraf.conf) to manage service configuration.
You can override the default configuration file with a VOLUME

```bash
docker run \
    -v /Users/gianarb/conf/telegraf.conf:/etc/telegraf.conf \
    -it gianarb/telegraf:0.10.4.1
```

Or you can use VOLUME and change the path with an env var

```bash
docker run \
    -e TELEGRAF_CONF=/Users/gianarb/conf/telegraf.conf \
    -v /Users/gianarb/git/conf/telegraf.conf \
    -it gianarb/telegraf:0.10.4.1
```
