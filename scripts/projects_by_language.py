import json

# Get list of usernames
username = 'glavin001'
with open('build/'+username+'/projects.json') as data_file:
    projects = json.load(data_file)

    langMap = {}
    for project in projects:
        langs = project['languages'].keys()
        for lang in langs:
            if lang in langMap:
                # Has Language already
                # Add project
                langMap[lang].append(project)
            else:
                langMap[lang] = [project]

    jsonStr = json.dumps(langMap)

    f = open('build/'+username+'/lang-projects.json','w')
    f.write(jsonStr) # python will convert \n to os.linesep
    f.close() # you can omit in most cases as the destructor will call it
