'''
import eventlet
pool = eventlet.GreenPool()
while True:
     pool.spawn(func,args )
'''

import eventlet
from eventlet.green import urllib2


urls = [
    "https://www.google.com/intl/en_ALL/images/logo.gif",
    "http://python.org/images/python-logo.gif",
    "http://us.i1.yimg.com/us.yimg.com/i/ww/beta/y3.gif",
]


def fetch(url):
    print "opening", url
    body = urllib2.urlopen(url).read()
    print "done with", url
    return url, body


pool = eventlet.GreenPool(200)
for url, body in pool.imap(fetch, urls):
    print "got body from", url, "of length", len(body)
