#!/bin/bash

/usr/share/elasticsearch/bin/elasticsearch -p /tmp/epid & /bin/bash /utils/wait-for-it.sh -t 0 localhost:9200 -- ./insert_script.sh;

wait $(cat /tmp/epid);
