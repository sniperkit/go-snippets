import json
import sys

user_id = {}
follows = []


def fix_line(line):
    return "{" + line[line.find("),")+2:]


for line in sys.stdin:
    try:
        fixed_line = fix_line(line)
        user = json.loads(fixed_line)
        user_id[user["login"]] = user["id"]
        follows.append((user["login"], user["follows"]))
    except Exception as e:
        print("Invalid line:", line, file=sys.stderr)
        print(e, file=sys.stderr)

follows_int = []

for (a, b) in follows:
    if a in user_id and b in user_id:
        print(user_id[a], user_id[b])
