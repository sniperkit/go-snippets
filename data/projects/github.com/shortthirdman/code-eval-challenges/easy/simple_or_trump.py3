import fileinput

v = {'2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8,
     '9': 9, '1': 10, 'J': 11, 'Q': 12, 'K': 13, 'A': 14}

def rank(c, t):
    r = v[c[0]]
    if c[-1] == t:
        r += 13
    return r

def win(l, r, t):
    lr, rr = rank(l, t), rank(r, t)
    if lr > rr:
        return l
    elif lr < rr:
        return r
    return l + ' ' + r

for line in fileinput.input():
    s = line.rstrip('\n').split()
    print(win(s[0], s[1], s[3]))
