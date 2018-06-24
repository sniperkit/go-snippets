#!/usr/bin/env python

import sys
from xml.etree import ElementTree
import json
from HTMLParser import HTMLParser
import feedparser

class MLStripper(HTMLParser):
    def __init__(self):
        self.reset()
        self.fed = []
    def handle_data(self, d):
        self.fed.append(d)
    def get_data(self):
        return ''.join(self.fed)

def extract_posts_from_rss(url):
    feed = feedparser.parse(url)
    for post in feed.entries:
        title = post.title
        link = post.link
        updated = post.updated
        summary = post.summary
#        print '###################################################'
        print '---------------------------------------------------'
        print title, updated
        print link
        s = MLStripper()
        s.feed(summary)
        print s.get_data()
#        print '---------------------------------------------------'
#        print content
#        print '---------------------------------------------------'


def extract_rss_urls_from_opml(filename):
    urls = []
    with open(filename, 'rt') as f:
        tree = ElementTree.parse(f)
    for node in tree.findall('.//outline'):
        url = node.attrib.get('xmlUrl')
        title = node.attrib.get('title')
        if url:
            urls.append( (title, url) )
    return urls

if __name__ == '__main__':
    if len(sys.argv) > 1:
        print("Parsing {}".format(sys.argv[1]))
        urls = extract_rss_urls_from_opml( sys.argv[1] )
        for url in urls:
            t, u = url
            print '###################################################'
            extract_posts_from_rss(u)
            print '###################################################'
