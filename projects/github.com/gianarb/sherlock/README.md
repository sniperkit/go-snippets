`sherlock` is a visualization tools for tracing that uses InfluxDB as backend.

More about how, what, why is available in the [`docs`](/docs/index.md) folder.

## Start daemon
```
make
./bin/scherlokd
curl -X GET -v localhost:8080/api/trace/34530
```

## quick start
We use `docker-compose`.

1. Checkout the project and go inside it

2. Docker compose up

```
docker-compose up
```
Builds influxdb and telegraf, ready to be used and well configured. Telegraf
with Zipkin plugin enable on port 9411. You can reach it from your host.
InfluxDB exposes 8096 to the host (not 8086).


3. Create InfluxDB database and Retention Policy

```
docker-compose exec influxdb influx
CREATE DATABASE traces
CREATE RETENTION POLICY weekly ON traces DURATION 1w REPLICATION 1
```
Traces can have a very high volume and cardinality, so we it is a good practice
go set a retention policy in order to keep the database reliable and usable.

4. Nothing more for now. You can point your applications instrumented with
   Opentracing sdk to 9411 to start storing traces in InfluxDB. [Or you can
   follow this
   tutorial](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/zipkin#example-input-trace)
