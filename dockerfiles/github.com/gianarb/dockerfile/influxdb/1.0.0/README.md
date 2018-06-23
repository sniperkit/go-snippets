# InfluxDB
[InfluxDB](http://influxdb.com) is an open-source, distributed, time series database with no external dependencies.

[hub.docker.com/gianarb/influxdb](https://hub.docker.com/r/gianarb/influxdb/)

## Run

```bash
docker run -it -p 8083:8083 -p 8086:8086 gianarb/influxdb
```

## Custom

The default configuration path is `/etc/influxdb_conf.toml` you can decide to use VOLUME to override it
```bash
docker run -it -p 8083:8083 -p 8086:8086 -v your/influxdb_conf.toml:/etc/influxdb_conf.toml gianarb/influxdb
```

Or you can change this patch with `$INFLUXDB_CONF`

```bash
docker run -it -e INFLUXDB_CONF=/influxdb.toml -v influxdb_conf.toml:/influxdb.toml gianarb/influxdb -config=/influxdb_conf.toml
```
