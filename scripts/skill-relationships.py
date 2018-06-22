import json

# Get list of usernames
with open('data/resume_usernames.json') as data_file:
    usernames = json.load(data_file)

    skillMap = {}

    # Get Resume for each username
    for index, username in enumerate(usernames):

        # Get Resume
        try:
            with open('data/resumes/'+username+'.json') as data_file:
                resume = json.load(data_file)
                skills = resume['skills']
                for skill in skills:
                    skillName = skill['name']
                    keywords = skill['keywords']
                    if skillName in skillMap:
                        skillMap[skillName].extend(keywords)
                    else:
                        skillMap[skillName] = keywords
        except Exception as e:
            print(e)

    #
    jsonStr = json.dumps(skillMap)

    f = open('data/skill-map.json','w')
    f.write(jsonStr) # python will convert \n to os.linesep
    f.close() # you can omit in most cases as the destructor will call it
