# Note

As of now I have:
- Swapped between keeping an index in RAM, or writing it to a more persistent storage
- Compute markov chains of sentences
- Tag the content (english text) for their POS (Part-Of-Speech) as defined by a brill tagger (see github.com/mvryan/fasttag)
- Train a markov chain on the penn treebank by Python's NLTK data (see http://www.nltk.org/nltk_data/ #18)
- Extract nouns and count their frequency in a document
- Pageranked some articles based on links
- Swapped the heading of the project from mysql, github.com/boltdb/bolt, to gzipping, zipping, flating, and pure text - also considered lzma, but that was not in the golang's std lib.
- Considered github.com/google/readahead for bulk RAM access of entire wikipedia, but that was too inefficient over indexing.

The heading now would be to use the IndexWiki for persistent storage (on a PC - 4GB RAM, > 100 GB drive, 2Ghz CPU) and read the indexes based on some score, term frequency needs implementation, of the links, references, titles, pages, or some other lexicon feature such as aiml.

Do what you want with it.

# WikiPagerankDB

## sample-boltdb.go
On Linux use the commands:

go get github.com/pressure679/WikiPagerankDB

go build $GOPATH/src/github.com/pressure679/WikiPagerankDB/sample-boltdb.go

cd $GOPATH/src/github.com/pressure679/WikiPagerankDB/

mkdir articles

cd articles

And depending on your machine (PC I recomment downloading the ~200MB-~2GB files, a 60+GB RAM/VRAM machine the single data dump is fine). Type the url https://dumps.wikimedia.org/enwiki/ into a webbrowser's URL input line, and a list of wikipedia data dump distributions should come up. Select the one you want.

wget https://dumps.wikimedia.org/enwiki/latest/enwiki-latest-pages-articles...

cd ..

./sample-boltdb

This is going to take some days, almost a week, with ~4GB RAM, ~2.4 GHz CPU, but do not worry about shutting down the PC, you can move the already used wikipedia dump files away from the directory and let the program continue from where it left off.

## sample-mysql.go
For this I recommend the mwdumper, see https://www.mediawiki.org/wiki/Manual:MWDumper
This I tested for ~20 minutes before I saw it only added ~50 articles, then I figured I used a base MySQL schema and decidede to make my own index of wikipedia and just use the raw xml files. (although MySQL or a custom compression or encoding would be nice).
For this install mysql and add a database named wikidb, then add your mysql username and password into the <username> and <password> field on line 71.

### Description of samples
sample-boltdb and sample-mysql run, and the same functionality of these are in tools.go but with the modification of making an index instead of the boltdb and/or mysql functionality.
