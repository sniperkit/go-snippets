import fileinput

for line in fileinput.input():
    c = s = 0
    while s < len(line):
        s = line.find("<--<<", s)
        if s == -1:
            break
        s += 4
        c += 1
    s = 0
    while s < len(line):
        s = line.find(">>-->", s)
        if s == -1:
            break
        s += 4
        c += 1
    print(c)
