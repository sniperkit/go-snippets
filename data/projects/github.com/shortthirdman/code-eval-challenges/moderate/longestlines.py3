import fileinput

n = 0
for line in (i.rstrip('\n') for i in fileinput.input()):
    if not n:
        n = int(line)
        c = [''] * n
        continue
    l = len(line)
    for ix, i in enumerate(c):
        if l > len(i):
            c.insert(ix, line)
            del c[-1]
            break
print(*c, sep='\n')
