import json
import os
import string

import numpy as np
import pandas as pd

# define constants

PRINTABLE = set(string.printable)
PATH_TO_DATA = os.getcwd()[:-7] + 'data/'

# Load skill-map.json

with open(PATH_TO_DATA + 'skill-map.json') as data_file:
	data = json.load(data_file)

# Generate parent and child lists

keys = data.keys()
parent_list = []
child_list = []

for key in data.keys():

	# Don't include non-english entries

	english = True

	for chr in key:

		if chr not in PRINTABLE:

			english = False
			break

	if not english:

		continue

	parent_list.append(key)

	children = []

	for child in data[key]:
		children.append(child.encode('utf-8'))


	child_list.append(children)


# get children into their own lists, along with related terms and parents

extracted_children = []
parents = []
related_terms = []
index = 0

for entry in child_list:

	for skill in entry:
		related = entry
		related.remove(skill) 
		extracted_children.append(skill)
		related_terms.append(related)
		parents.append(parent_list[index])


	index += 1

# Generate dataframe skill_rel

skill_rel = pd.DataFrame({"keywords": pd.Series(extracted_children),
						  "parents": pd.Series(parents),
						  "related_terms": pd.Series(related_terms)})


# Hand crafted sorting

# Move all rows with SQL to Databases
SQL_parents_index = np.ndarray.tolist(skill_rel[~skill_rel['parents'].str.contains("^((?!SQL).)*$")].index.values)
skill_rel['parents'].ix[SQL_parents_index] = 'Databases'

# Move all rows with database to Databases
database_parents_index = np.ndarray.tolist(skill_rel[~skill_rel['parents'].str.contains("^((?!database).)*$")].index.values)
skill_rel['parents'].ix[database_parents_index] = 'Databases'

# Move all rows with Database to Databases
Database_parents_index = np.ndarray.tolist(skill_rel[~skill_rel['parents'].str.contains("^((?!Database).)*$")].index.values)
skill_rel['parents'].ix[Database_parents_index] = 'Databases'

# Move Backend Development to Backend
backend_parents_index = np.ndarray.tolist(skill_rel[~skill_rel['parents'].str.contains("^((?!Backend).)*$")].index.values)
skill_rel['parents'].ix[backend_parents_index] = 'Backend'

# Move Frontend Development to Frontend
frontend_parents_index = np.ndarray.tolist(skill_rel[~skill_rel['parents'].str.contains("^((?!Front).)*$")].index.values)
skill_rel['parents'].ix[frontend_parents_index] = 'Frontend'

# Move Languages to Programming
Languages_parents_index = np.ndarray.tolist(skill_rel[~skill_rel['parents'].str.contains("^((?!Languages).)*$")].index.values)
skill_rel['parents'].ix[Languages_parents_index] = 'Programming'

# Move Software to Development
Software_parents_index = np.ndarray.tolist(skill_rel[~skill_rel['parents'].str.contains("^((?!Software).)*$")].index.values)
skill_rel['parents'].ix[Software_parents_index] = 'Development'

# Move DevOps to Development
DevOps_parents_index = np.ndarray.tolist(skill_rel[~skill_rel['parents'].str.contains("^((?!DevOps).)*$")].index.values)
skill_rel['parents'].ix[DevOps_parents_index] = 'Development'


# print skill_rel['keywords'].value_counts()
# print skill_rel['parents'].value_counts()

# print skill_rel

# export unique wordlist

pd.Series(skill_rel["keywords"].unique()).to_json(PATH_TO_DATA + "unique_keywords.json",orient='records')






