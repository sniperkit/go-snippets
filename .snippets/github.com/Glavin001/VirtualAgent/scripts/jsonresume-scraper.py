# Example: http://registry.jsonresume.org/thomasdavis.json

import json
import urllib2

# Get list of usernames
with open('data/resume_usernames.json') as data_file:
    usernames = json.load(data_file)
    # Request resume for username

    for index,username in enumerate(usernames):
        # Request resume
        resumeUrl = 'http://registry.jsonresume.org/'+username+'.json'
        print(index, resumeUrl)
        req = urllib2.Request(resumeUrl)
        # req.add_header('Content-Type', 'application/json; charset=utf-8')
        req.add_header('Accept', 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8')
        req.add_header('Host', 'registry.jsonresume.org')
        req.add_header('User-Agent','Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.86 Safari/537.36')
        '''
        Accept:text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
        Accept-Encoding:gzip, deflate, sdch
        Accept-Language:en-US,en;q=0.8
        Cache-Control:no-cache
        Connection:keep-alive
        Cookie:__cfduid=d449297d8b9bbe53f972a3196932420e61461968967; connect.sid=s%3AU7ipVBaJHbmXk4EfzQ9Rdf0QxsPxy20O.GhQvjGfOjrmXCxOVViZRL%2F2m4LVu8oPsDEhzPAJGS%2FY; _ga=GA1.2.1843410159.1461968971
        Host:registry.jsonresume.org
        Pragma:no-cache
        Upgrade-Insecure-Requests:1
        User-Agent:Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.86 Safari/537.36
        '''
        try:
            resp = urllib2.urlopen(req)
            resume = resp.read()

            f = open('data/resumes/'+username+'.json','w')
            f.write(resume) # python will convert \n to os.linesep
            f.close() # you can omit in most cases as the destructor will call it
        except Exception as e:
            print(e)

