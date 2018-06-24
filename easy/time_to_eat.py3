import fileinput

for line in fileinput.input():
    print(*sorted(line.split(), reverse=True))
