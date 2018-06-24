import fileinput

for line in fileinput.input():
    s = line.split()
    c, o, r, v = 0, 1, 0, 0
    for i in s[1]:
        if i in 'abcdefghijklmnopqrstuvwxyz':
            v, c = v*10 + int(s[0][c]), c + 1
        elif i == '+':
            o, r, v = 1, r + v*o, 0
        elif i == '-':
            o, r, v = -1, r + v*o, 0
    print(r + v*o)
