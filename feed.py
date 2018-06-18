#!/usr/bin/env python

import sys
import json
import feedparser
from xml.etree import ElementTree

def extract_posts_from_rss(url):
    print("Parsing {}".format( url ))
    feed = feedparser.parse(url)
    for post in feed.entries:
        title = post.title
        link = post.link
        updated = post.updated
        summary = post.summary
        content = post.content
        print '###################################################'
        print '---------------------------------------------------'
        print title, updated
        print link
        print '---------------------------------------------------'
        print summary
        print '---------------------------------------------------'
        print content
        print '---------------------------------------------------'

if __name__ == '__main__':
    url = 'http://googleappengine.blogspot.com/atom.xml'
    if len(sys.argv) > 1:
        url = sys.argv[1]
    extract_posts_from_rss( url )
