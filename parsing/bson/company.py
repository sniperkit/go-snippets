import sys


def read_issues(source, outdir):
    pairs = []
    for line in source:
        try:
            ss = line.strip().split(",")
            if len(ss) == 1:
                user, comp = ss[0], ""
            elif ss[1].startswith("None"):
                user, comp = ss[0], ""
            else:
                user, comp = ss[0], ss[1]
            pairs.append((user, comp))
        except Exception as e:
            print("Problem reading line:", line.strip(), file=sys.stderr)
            print(e, file=sys.stderr)

    # Write
    f = open(outdir + ".txt", "w")
    for pair in pairs:
        s = "%s,%s\n" % (pair[0], pair[1])
        f.write(s)
        f.flush()
    f.close()

if __name__ == "__main__":
    if(len(sys.argv) < 2):
        print("Please specify an output directory")
    elif(len(sys.argv) == 2):
        read_issues(sys.stdin, sys.argv[1])
    else:
        print(sys.argv[1], sys.argv[2])
        read_issues(open(sys.argv[1]), sys.argv[2])
