# Chronograf
[Chronograf](https://influxdata.com/time-series-platform/chronograf/) is a
simple to install graphing and visualization application that you deploy behind
your firewall to perform ad hoc exploration of your InfluxDB data, It includes
support for templates and a library of intelligent, pre-configured dashboards
for common data sets.

## Start
```
docker pull gianarb/chronograf:0.10.0
docker run -p 10000:10000 -v /path/database:/data gianarb/chronograf:0.10.0
```
Open your browser and you will se this application run on port `10000` `-v
/path/database:/data` helps to manage the database store with your dashboard
and widget.

## Custom
```
docker pull gianarb/chronograf:0.10.0
docker run -p 10000:10000 -v /path/your/conf.toml:/chronograf.toml -v /path/database:/data gianarb/chronograf:0.10.0
```
In this way you can use your configuration file.
