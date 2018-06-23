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

jsonObject = json.load(open(p + "\data\\jobs.json", 'r'))

tags = set()
for dictionary in jsonObject['jobs']:
    for key in dictionary:
        if key == 'name':
            tags.add(dictionary[key])

#print tags
