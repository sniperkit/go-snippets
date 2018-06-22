import sys


def read_issues(collabs_tmp, users_tmp, repos_tmp, outdir):
    collab_dict = dict()
    for line in collabs_tmp:
        ss = line.strip().split(",")
        val = "%s/%s" % (ss[1], ss[2])
        key = ss[0]
        if key in collab_dict.keys():
            collab_dict[key].append(val)
        else:
            collab_dict[key] = [val]
    
    repo_dict = dict()
    for line in repos_tmp:
        ss = line.strip().split(",")
        key = "%s/%s" % (ss[1], ss[0])
        val = ss[2]
        if key in repo_dict.keys():
            repo_dict[key].append(val)
        else:
            repo_dict[key] = [val]
    
    fout = open(outdir + ".txt", "w")
    user_dict = dict()
    for line in users_tmp:
        ss = line.strip().split(",")
        user = ss[0]
        if user in collab_dict.keys():
            for repo in collab_dict[user]:
                if repo in repo_dict.keys():
                    for lang in repo_dict[repo]:
                        s = "%s,%s,%s\n" % (user,repo,lang)
                        fout.write(s)
                else:
                    s = "%s,%s,<UNKNOWN>\n" % (user,repo)
                    fout.write(s)
        fout.flush()
    fout.close()
    

if __name__ == "__main__":
    if(len(sys.argv) < 2):
        print("Please specify an output directory")
    elif(len(sys.argv) == 2):
        read_issues(sys.stdin, sys.argv[1])
    else:
        print(sys.argv[1], sys.argv[2], sys.argv[3], sys.argv[4])
        read_issues(open(sys.argv[1]), open(sys.argv[2]), open(sys.argv[3]), sys.argv[4])
