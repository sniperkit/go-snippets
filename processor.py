import re
from bs4 import BeautifulSoup
    
def visible(element):
    if element.parent.name in ['style', 'script', '[document]', 'head', 'title']:
        return False
    elif re.match('<!--.*-->', str(element.encode('utf-8'))):
        return False
    return True

def trim(element):
    if element.isspace():
        return False
    return True
    
def inside_trim(element):
    return "".join(element.split()) + " "


def process(soup):
    data = soup.findAll(text=True)
    result = filter(visible, data)
    result = filter(trim, result)
    result = map(inside_trim, result)
    return "".join(result).rstrip()
