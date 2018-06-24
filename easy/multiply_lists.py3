import fileinput

for line in fileinput.input():
    s0, s1 = [i.split() for i in line.split('|')]
    print(*[str(int(i) * int(s1[ix])) for ix, i in enumerate(s0)])
