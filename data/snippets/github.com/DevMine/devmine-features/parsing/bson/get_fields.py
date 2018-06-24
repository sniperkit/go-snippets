from bsonstream import KeyValueBSONInput
import sys


def read_field(record, f):
    r = record
    for p in f.split("/"):
        r = r[p]

    return str(r)


def read_file(source, output, fields):
    for _, record in KeyValueBSONInput(fh=source):
        try:
            print >> output, ",".join([read_field(record, f) for f in fields])
        except KeyError:
            pass
        except Exception as e:
            print >> sys.stderr, e


if __name__ == "__main__":
    if sys.argv[1] == "--stdin":
        infile = sys.stdin
    else:
        infile = open(sys.argv[1], "rb")

    if sys.argv[2] == "--stdout":
        outfile= sys.stdout
    else:
        outfile = open(sys.argv[2], "w")

    read_file(infile, outfile, sys.argv[3:])
