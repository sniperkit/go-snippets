import fileinput

for line in fileinput.input():
    if not line == '\n':
        print(*reversed(line.split()))
