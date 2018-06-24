import fileinput

for line in fileinput.input():
    r, s = 0, [int(i) for i in line.translate({ord('x'):ord(' '), ord('|'):None}).split()]
    while s[1] != s[3]:
        r, s[0], s[1], s[2], s[3] = r+s[0], s[1]-1, s[0], s[1]-s[3], s[2]
    print(r+s[2])
