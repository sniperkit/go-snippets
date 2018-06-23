"""
Author: Conor Scott
Date: 04/29/16
Volta Hackathon

Purpose: Scraping StackOverflow for job postings

Requirements: feedparser
"""

import feedparser
import os

#Initializes feed to a RSS link
def initializeFeed(link):
	#Link to RSS Feed
	#Initialize RSS feed
	feed = feedparser.parse(link)
	return feed

#Get unique attributes of a type, e.g. 'title' or 'category' as a valid key 
def getUniqueAttributes(Feed, key, allowDuplicates = True):
	uniqueAttributes = None
	if allowDuplicates:
		uniqueAttributes = []
	else:
		uniqueAttributes = set()


	for post in Feed.entries:

		if 'category' not in post:
			continue
		if 'location' not in post:
			continue
		elif key in post:
			if not allowDuplicates:
				uniqueAttributes.add(post[key.encode('utf-8')])
			else:
				uniqueAttributes.append(post[key.encode('utf-8')])
		#else:
		#	if key == 'location':
		#		uniqueAttributes.append('N/A')
		#	if key == 'category':
		#		uniqueAttributes.append('N/A')
	return uniqueAttributes

#Initialize feed
myFeed = initializeFeed('http://stackoverflow.com/jobs/feed')

"""
Compute all unique categories of stackoverflow job postings, all 
relevant categories are in respective sets avoiding duplicates
"""

categories = list(getUniqueAttributes(myFeed, 'category'))
#print len(categories)
titles = list(getUniqueAttributes(myFeed, 'title'))

links = list(getUniqueAttributes(myFeed, 'link'))

locations = list(getUniqueAttributes(myFeed, 'location'))

dates = list(getUniqueAttributes(myFeed, 'date'))

authors = list(getUniqueAttributes(myFeed, 'author'))

updated = list(getUniqueAttributes(myFeed, 'updated'))

descriptions = list(getUniqueAttributes(myFeed, 'description'))

for des in descriptions:
	pass

"""
Author: Conor Scott
Date: 04/29/16
Volta Hackathon

Purpose: Unique tags (Like Java) from stackoverflow are obtained from 
the json file I copied from their website API that obtains 
them 100 at a time, and compiled 1000 of the most popular
tags

Requirements: just python (comes with json)
"""
import json, os

p = os.path.abspath('..')
jsonObject = json.load(open(p + "/data/tags.json", 'r'))
tags = set()
for dictionary in jsonObject['items']:
    for key in dictionary:
        if key == 'name':
            tags.add(dictionary[key])

def loadResume(filename):
	p = os.path.abspath('..')
	jsonObject = json.load(open(p + filename, 'r'))
	return jsonObject

def getSkills(resume):
	return resume['skills']

def getUniqueSkills(skills):
	a = set()
	for keyword in skills: 
		for item in keyword['keywords']:
			a.add(item)
	return a

allCategories = set(categories + list(tags))

actualCategories = []

for index, description in enumerate(descriptions):
	currentCategories = []
	for category in allCategories:
		if description.find(category) != -1:
			currentCategories.append(category)		
	actualCategories.append(currentCategories)

#print len(categories) = 949
#print len(tags) = 1000
#print len(allCategories) = 1110
#print allCategories
#print allCategories

#print allCategories


#Get Aaron's Skills from his resume
#uniqueSkills = getUniqueSkills(getSkills(loadResume("\data\\resumes\\aaron.json")))

#unionOfCategoriesAndSkills = list(set(categories) | set(uniqueSkills))

uniques = sorted(list(set(categories) | set(tags)))

jobpostings = [titles, descriptions, dates, locations, authors, links]
jobpostingslengths = [len(jobpostings[i]) for i in range(0, len(jobpostings))]
jobkeys = ["title", "description", "date_created", "full_time", "location", "company_name", "source", "apply", "post_url"]


nonUTF8 = []
p = os.path.abspath('..')
with open(p + "\data\\stackoverflow_unique_values.json", 'w') as f:
	f.write("[")
	for index, tag in enumerate(uniques):
		a = tag.decode('utf-8')
		if index == len(uniques)-1:
			f.write(a + "]")
		else:
			f.write('"' + a + '"' + ", ")
	f.close()

mainJobMap = {}
jobMap = {}
job_list = []

for i in range(0, len(titles)):
	
	s = {}

	s['title'] = jobpostings[0][i].encode('utf-8')
	s['description'] = jobpostings[1][i].encode('utf-8') + os.linesep.encode('utf-8' ) + os.linesep.encode('utf-8') + 'This job is looking for: '.encode('utf-8') + ', '.encode('utf-8').join(actualCategories[i]).encode('utf-8')
	s['date_created'] = jobpostings[2][i].encode('utf-8')
	s['full_time'] = "N/A"
	s['location'] = jobpostings[3][i].encode('utf-8')
	s['company_name'] = jobpostings[4][i].encode('utf-8')
	s['company_url'] = "N/A"
	s['source'] = "StackOverflow Job Postings"
	s['apply'] = "how to contact and apply"
	s['post_url'] = jobpostings[5][i].encode('utf-8')

	job_list.append(s)

#mainJobMap['jobs'] = jobMap
#jsonString = json.dumps(job_list)
#f.write(jsonString)

PATH_TO_DATA = os.getcwd()[:-7] + 'data/'

with open(PATH_TO_DATA + 'StackOverflow_Jobs.json', 'w') as outfile:
	json.dump(job_list, outfile)
