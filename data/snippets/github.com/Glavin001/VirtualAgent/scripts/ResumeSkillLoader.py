import json
import os

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

#Get Aaron's Skills from his resume
getUniqueSkills(getSkills(loadResume("\data\\resumes\\aaron.json")))