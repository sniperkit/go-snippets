from elasticsearch import Elasticsearch
import urllib2
import requests
import json
import os

es = Elasticsearch([{'host': 'localhost', 'port': 9200}])

PATH_TO_DATA = os.getcwd()[:-7] + 'data/'

with open(PATH_TO_DATA + 'unique_keywords.json') as data_file:
	data = json.load(data_file)

skill_list = []

for skill in data:
	skill_list.append(skill)

with open(PATH_TO_DATA + 'GitHub_Jobs.json') as data_file:
	data = json.load(data_file)


for i in range(len(data)):

	es.index(index='github_jobs', doc_type='test-type',id = i, body = data[i])	

print es.search(index='github_jobs', body={"query": {"match": {'description':'interest'}}})

from datetime import datetime
from elasticsearch_dsl import DocType, String, Date, Integer
from elasticsearch_dsl.connections import connections

# Define a default Elasticsearch client
connections.create_connection(hosts=['localhost'])

class Article(DocType):
    title = String(analyzer='snowball', fields={'raw': String(index='not_analyzed')})
    body = String(analyzer='snowball')
    tags = String(index='not_analyzed')
    published_from = Date()
    lines = Integer()

    class Meta:
        index = 'blog'

    def save(self, ** kwargs):
        self.lines = len(self.body.split())
        return super(Article, self).save(** kwargs)

    def is_published(self):
        return datetime.now() > self.published_from

# create the mappings in elasticsearch
Article.init()

# create and save and article
article = Article(meta={'id': 42}, title='Hello world!', tags=['test'])
article.body = ''' looong text '''
article.published_from = datetime.now()
article.save()

article = Article.get(id=42)
print(article.is_published())

# Display cluster health
print(connections.get_connection().cluster.health())

