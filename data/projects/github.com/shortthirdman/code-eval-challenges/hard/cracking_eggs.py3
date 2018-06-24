import fileinput

ht = {}
def floors(e, s):
    if e == 0 or s == 0:
        return 0
    elif e == 1:
        return s
    elif s == 1:
        return 1
    elif e > s:
        return floors(s, s)
    if (e, s) not in ht:
        ht[(e, s)] = floors(e-1, s-1) + floors(e, s-1) + 1
    return ht[(e, s)]

for line in fileinput.input():
    e, s = [int(i) for i in line.split()]
    n = 0
    while floors(e, n) < s:
        n += 1
    print(n)
