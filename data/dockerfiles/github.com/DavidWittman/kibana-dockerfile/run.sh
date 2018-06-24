#!/bin/sh

ELASTICSEARCH_URL=${ELASTICSEARCH_URL:-http://localhost:9200}

exec /opt/kibana/bin/kibana -e ${ELASTICSEARCH_URL} "$@"
