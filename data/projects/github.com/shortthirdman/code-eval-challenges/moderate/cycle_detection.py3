import fileinput

for line in fileinput.input():
    s = line.split()
    f, j = [], {}
    while s:
        i = s.pop(0)
        if i not in j:
            if not s:
                break
            j[i] = s[0]
            if f:
                break
        else:
            if i in f:
                break
            f.append(i)
    print(*f)
