#!/bin/bash
# Create indices here!

echo Creating movie index...

curl -X PUT "http://localhost:9200/movies/" -H "Content-Type: application/json" -d '{
    "mappings": {
        "movie": {
            "properties": {
                "title": {
                    "type": "text",
                    "fields": {
                        "completion": {
                            "type": "completion"
                        }
                    }
                },
                "image": {
                    "type": "text"
                },
                "imdb_meta": {
                    "type": "nested",
                    "properties": {
                        "genre": {
                            "type": "text"
                        },
                        "mpaa_rating": {
                            "type": "text"
                        },
                        "score": {
                            "type": "float"
                        }
                    }
                },
                "rotten_tomatoes_meta": {
                    "type": "nested",
                    "properties": {
                        "tomato_score": {
                            "type": "integer"
                        },
                        "popcorn_score": {
                            "type": "integer"
                        },
                        "theater_release_date": {
                            "type": "text"
                        },
                        "mpaa_rating": {
                            "type": "text"
                        },
                        "synopsis": {
                            "type": "text"
                        },
                        "synonpsis_type": {
                            "type": "text"
                        },
                        "runtime": {
                            "type": "text"
                        }
                    }
                }
            }
        }
    }
}';

echo Finished creating movie index