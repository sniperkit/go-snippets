from __future__ import print_function
import sys
import os
import tarfile
from bsonstream import KeyValueBSONInput


def read_languages(sourcedir, outdir):
    # Read
    user_langs = {}
    languages = set()

    for fname in os.listdir(sourcedir):
        print("Opening " + fname)
        f = tarfile.open(sourcedir + "/" + fname).extractfile("dump/github/repos.bson")
        print("Parsing " + fname)
        stream = KeyValueBSONInput(fh = f)

        for _, repo in stream:
            try:
                user = repo["owner"]["login"]
                language = repo["language"]
                size = repo["size"]
                if language:
                    languages.add(language)
                    if user in user_langs:
                        if language not in user_langs[user]:
                            user_langs[user][language] = 0
                        user_langs[user][language] += int(size)
                    else:
                        user_langs[user] = {language: int(size)}
            except Exception as e:
                print("Problem reading line:", line.strip(), file=sys.stderr)
                print(e, file=sys.stderr)

        f.close()

    # Write
    print("Writing to files")
    files = {}
    for language in languages:
        files[language] = open(outdir + "/" + language.replace("/", "|"), "w")

    for user in user_langs:
        for l in user_langs[user]:
            if(user_langs[user][l] > 0):
                files[l].write(user + ',' + str(user_langs[user][l]) + '\n')

    for f in files.values():
        f.flush()
        f.close()

if __name__ == "__main__":
    if(len(sys.argv) < 3):
        print("Usage: python languages.py inputdir outputdir")
    else:
        read_languages(sys.argv[1], sys.argv[2])
