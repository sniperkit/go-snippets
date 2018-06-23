#!/bin/bash

curl -X PUT localhost:9200/iloveindexes -d \
'{
    "mappings" : {
        "index_of" : {
            "properties" : {
                "host" : {
                    "type" : "string",
                    "index" : "not_analyzed" 
                }
            }
        }
    }
}'
