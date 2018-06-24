"""print a 2D array (n x m) in spiral order (clockwise)"""
import fileinput

for line in fileinput.input():
    st = line.rstrip('\n').split(';')
    m = st[2].split(' ')
    r = []

    a = int(st[0])
    b = int(st[1])
    c = la = lb = 0
    ha = a - 1
    hb = b - 1

    while True:
        while c % b < hb:
            r.append(m[c])
            m[c] = -1
            c += 1
        if (c + b >= a * b) or (m[c + b] == -1):
            break
        la += 1
        while int(c / b) < ha:
            r.append(m[c])
            m[c] = -1
            c += b
        if (c % b == 0) or (m[c - 1] == -1):
            break
        hb -= 1
        while c % b > lb:
            r.append(m[c])
            m[c] = -1
            c -= 1
        if m[c - b] == -1:
            break
        ha -= 1
        while int(c / b) > la:
            r.append(m[c])
            m[c] = -1
            c -= b
        if m[c + 1] == -1:
            break
        lb += 1

    r.append(m[c])
    print(' '.join(r))
