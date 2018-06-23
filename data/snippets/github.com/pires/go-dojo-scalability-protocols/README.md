# go-dojo-scalability-protocols
Go experiment with [Scalability Protocols](http://bravenewgeek.com/a-look-at-nanomsg-and-scalability-protocols/).

## Pre-requisites

* `go` 1.4.2+
* [`gb`](http://getgb.io)

## Build

```
gb build all
```

## Run

For all available flags please run
```
./bin/sp --help
```

### Locally

```
./bin/sp -mode=receiver
```

#### Sender

To send 1 million messages:
```
./bin/sp -mode=sender -num_messages=1000000
```
To terminate, press _CTRL^C_.

### Over LAN

Assuming `receiver` instance is running on host `192.168.0.123`.

#### Receiver

To terminate, press _CTRL^C_.
```
./bin/sp -mode=receiver -address=192.168.0.123
```

#### Sender

To send 1 million messages:
```
./bin/sp -mode=sender -address=192.168.0.123 -num_messages=1000000
```
To terminate, press _CTRL^C_.
