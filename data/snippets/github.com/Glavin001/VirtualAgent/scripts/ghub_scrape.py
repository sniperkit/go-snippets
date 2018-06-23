import json
import urllib2
import os
import string

PATH_TO_DATA = os.getcwd()[:-7] + 'data/'

def get_jobs(page=0):
	''' Returns all jobs on a specified github careers page as a json object.
	'''

	response = urllib2.urlopen('https://jobs.github.com/positions.json' + '?page=' + str(page))
	json_object = json.load(response)
	for job in json_object:
		job.pop("company_logo", None)
		job.pop("id", None)
		if job['type'] == 'Full Time':
			job['full-time'] = True
		job.pop("type", None)
		job['company_name'] = job['company']
		job.pop("company", None)
		job['data_created'] = job['created_at']
		job.pop("created_at", None)
		job['post_url'] = job['url']
		job.pop("url", None)
		job['source'] = 'GitHub Careers'
		job['apply'] = job['how_to_apply']
		job.pop("how_to_apply", None)

	return json_object

def get_all_jobs():

	job_list = []
	index = 0
	while get_jobs(index) != []:

		job_list.extend(get_jobs(index))
		index += 1

	return job_list

with open(PATH_TO_DATA + 'GitHub_Jobs.json', 'w') as outfile:

    json.dump(get_all_jobs(), outfile)








