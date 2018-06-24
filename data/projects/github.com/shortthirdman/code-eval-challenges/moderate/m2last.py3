import fileinput

for line in fileinput.input():
    s = line.split()
    n = len(s) - 1 - int(s[-1])
    if n >= 0:
        print(s[n])
