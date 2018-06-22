# Goloso

_Goloso is an Ansible callback server over NSQ_

The idea here it to have an automatic bootstrap mechanism using a callback system to run Ansible playbooks remotely.

Whenever a instance goes up, it posts a message saying “hi guys, I’m alive” to a central message broker. Them comes **Goloso**, our engine which consumes from a “bootstrap” topic and orders the playbook run into that fresh new and sweet instance.

This goloso points to a single Ansible playbook, the master entrypoint. This playbook knows what to choose based on the variables/tags received.


## System Overview

The system has three main components.

* The message producer or new instance which will post a JSON message to NSQ.

* NSQ, the message broker.

* Goloso, consuming messages from NSQ and running Ansible playbooks.


## Install Goloso

    go get github.com/gophergala/goloso

## Goloso Dependencies

    go get github.com/bitly/go-nsq
    go get github.com/bitly/nsq/util
    go get github.com/boltdb/bolt

## Install NSQ

    brew install nsq
    
## Run NSQ

    nsqlookupd
    nsqd --lookupd-tcp-address=127.0.0.1:4160
    nsqadmin --lookupd-http-address=127.0.0.1:4161
    open http://localhost:4171
    
    
## cloud-init

This shell script `scripts/cloud_init_bootstrap.sh` is the spark of the system.  This script should be set under your instance user-data and It will produce the message above.

    {
        "event": "bootstrap",
        "uuid": "51DC9184-4499-4FDA-9EA3-4D19CD07486A",
        "instance_id": "i-806086",
        "ipaddress": "127.0.0.7",
        "os": "unknown"
    }

## Install Ansible

    brew install ansible

## Future

 * web admin interface
 * replay events
