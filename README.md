# Kibana Dockerfile

A minimal Docker repository for Kibana 4.x based on [Alpine Linux](https://www.alpinelinux.org).

## Installation

``` bash
$ docker pull objectrocket/kibana
```

## Usage

``` bash
$ docker run -d -e "ELASTICSEARCH_URL=https://user:pass@es.example.com:9200" -p 5601:5601 objectrocket/kibana
```

Then you can visit your Docker host at port 5601 to view the Kibana dashboard.
