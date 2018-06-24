#/usr/bin/env python

from sklearn.feature_extraction.text import TfidfVectorizer
from sys import argv
from os import walk,write, mkdir, listdir
from os import path

NUMB_OF_FEATURES=30

if len(argv) < 2:
    print "Please provide a corpus file/directory"
    exit(0)


def read_corpus_from_file(corpus_file):
    corpus = []
    scores= []
    f = open(corpus_file, 'r')
    while True:
        lines=f.readline()
        lines=lines.strip()
        lines=lines.split(",")
        if len(lines)>0 and lines[0]!='':
            corpus.append(lines[0])
            scores.append(lines[1])
        else:
            break
    return corpus,scores

def read_corpus(corpus_file):
    corpus = []
    scores=[]
    filenames =[]
    counts = []
    if not path.isdir(corpus_file):
        corpus, scores=read_corpus_from_file(corpus_file)
        filenames=[corpus_file]
        counts=[len(corpus)]
    else:
        for f in listdir(corpus_file):
            c, s = read_corpus_from_file(corpus_file+"/"+f)
            corpus.extend(c)
            scores.extend(s)
            filenames.append(f)
            counts.append(len(c))
    return corpus, scores, filenames, counts
        
def tfidf(corpus):
    vectorizer = TfidfVectorizer(stop_words='english', lowercase=True, max_features=NUMB_OF_FEATURES)
    return vectorizer.fit_transform(corpus).toarray()


def create_results(matrix, scores, outfile, filenames, counts):
    count = 0
    mkdir(outfile)
    offset=0
    for idx in range(0, len(filenames)):
        out = open(outfile+"/"+filenames[idx], 'w')
        for i in range(1,NUMB_OF_FEATURES+1):
            out.write("feat-%d," % i)
        out.write("class\n")
        for i in range(offset, offset+counts[idx]):
            all_zeros=True
            res_out=""
            for j in range(0, len(matrix[i])):
                res_out+="%.5f," % (matrix[i][j])
                all_zeros = all_zeros and (matrix[i][j]==0.0)
            res_out+=str(scores[i])+"\n"
            if all_zeros:
                count+=1
            else:
                out.write(res_out)
        offset+=counts[idx]
        print "Found %d elements with all zero elements" % count
        out.close()



corpus_file=argv[1]
output_file=argv[2]

corpus,scores,filenames,counts=read_corpus(corpus_file)
res = tfidf(corpus)
create_results(res, scores,output_file, filenames, counts)
