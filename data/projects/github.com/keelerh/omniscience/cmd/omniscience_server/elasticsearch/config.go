package elasticsearch

const (
	Index   = "omniscience"
	Mapping = `
		{
			"settings":{
				"number_of_shards": 1,
				"number_of_replicas": 0
			},
			"mappings":{
				"_doc":{
					"properties":{
						"id":{
							"type":"keyword"
						},
						"title":{
							"type":"text",
							"store": true
						},
						"description":{
							"type":"text",
							"store": true
						},
						"service":{
							"type":"text",
							"store": true
						},
						"content":{
							"type":"text"
						},
						"url":{
							"type":"text",
							"store": true
						},
						"last_modified":{
							"type":"date",
							"store": true
						},
						"suggest_field":{
							"type":"completion"
						}
					}
				}
			}
		}`
)
