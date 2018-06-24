import gevent
from gevent import pool
from gevent import monkey
monkey.patch_all()
import requests
from warcio.archiveiterator import ArchiveIterator
import zmq
import ujson
from bs4 import BeautifulSoup
import processor
import operator

context = zmq.Context()
socket = context.socket(zmq.PUB)
socket.bind("tcp://0.0.0.0:5555")

POOL = pool.Pool(100)

def word_count(string):
    counts = dict()
    words = string.split()

    for word in words:
        if word in counts:
            counts[word] += 1
        else:
            counts[word] = 1

    return counts

def print_records(url):
    url = url.strip()
    resp = requests.get(url, stream=True)
    for record in ArchiveIterator(resp.raw):
        if record.rec_type == 'warcinfo':
            pass

        elif record.rec_type == 'response':
            # print(record.rec_headers)
            if not record.http_headers:
                continue
            if record.http_headers.get_header('Content-Type') == 'text/html':
                try: 
                    soup = BeautifulSoup(record.content_stream().read().decode("utf-8"))
                except Exception as e:
                    # print(e)
                    continue


                # Process record here, maybe spacy

                text = processor.process(soup)


                counts = word_count(text)

                top_3_words = [x[0] for x in sorted(counts.items(), key=operator.itemgetter(1), reverse=True)[:2]]

                node = record.rec_headers.get_header('WARC-Target-URI')

                outlinks = ",".join([link['href'] for link in soup.find_all('a', href=True)])

                msg = bytes(ujson.dumps({"Node":node,"Keywords":",".join(top_3_words), "Outlinks":outlinks, "Score":1.0}), "utf-8")

                socket.send(msg)
                # print(msg.decode("utf-8"))
    

with open("warc copy.txt", "r") as textfile:
    urls = textfile.readlines()
for url in urls:
    POOL.spawn(print_records, "https://commoncrawl.s3.amazonaws.com/"+urls[0])
POOL.join()
socket.send(bytes("done", "utf-8"))

# print_records('https://archive.org/download/ExampleArcAndWarcFiles/IAH-20080430204825-00000-blackbook.warc.gz')