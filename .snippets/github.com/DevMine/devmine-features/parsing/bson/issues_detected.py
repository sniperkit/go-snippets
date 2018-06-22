import sys


def read_issues(source, outdir):
    users_count = dict()
    for line in source:
        try:
            user = line.strip()[:-1]
            user = str(user)
            if user in users_count:
                users_count[user] = users_count[user] + 1
            else:
                users_count[user] = 1
        except Exception as e:
            print("Problem reading line:", line.strip(), file=sys.stderr)
            print(e, file=sys.stderr)

    # Write
    f = open(outdir + ".txt", "w")

    for user in users_count.keys():
        s = '%s,%d\n' % (str(user), users_count[user])
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
