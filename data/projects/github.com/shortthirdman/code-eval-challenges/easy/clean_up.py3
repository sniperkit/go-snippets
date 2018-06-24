import fileinput, re

for line in fileinput.input():
    print(*re.findall("[a-z]+", line.lower()))
