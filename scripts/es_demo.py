import requests
import os
import json


# Constants

PATH_TO_DATA = os.path.abspath('data')+'/'

# Determine if ES is operational

res = requests.get('http://localhost:9200')
print res.content

# Connect to node

from elasticsearch import Elasticsearch
es = Elasticsearch([{'host': 'localhost', 'port': 9200}])

# Import and Combine data

with open(PATH_TO_DATA + 'GitHub_Jobs.json') as data_file:
	GH_data = json.load(data_file)

with open(PATH_TO_DATA + 'StackOverflow_Jobs.json') as data_file:
	SO_data = json.load(data_file)

combined_jobs = GH_data + SO_data

# Export Combined data for later

with open(PATH_TO_DATA + 'Combined_Jobs.json', 'w') as outfile:

    json.dump(combined_jobs, outfile)

# Index Data

# Only need to run indexing code ONCE, uncomment it if needed: be warned 66 Second run time

# for i in range(len(combined_jobs)):
#
# 	es.index(index='combined_jobs', doc_type='test-type',id = i, body = combined_jobs[i])

# Example Query + returns highest scoring job's post formated as specified

print combined_jobs[int(es.search(index='combined_jobs', body ={"query": {"match" : { "description": "C++ AND Java"}}})['hits']['hits'][0]['_id'].encode('utf-8'))]







