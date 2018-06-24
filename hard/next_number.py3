import fileinput

def dsig(a):
    r = 0
    while a > 0:
        if a % 10 > 0:
            r += 1 << (3 * (a % 10))
        a //= 10
    return r

for line in fileinput.input():
    d = int(line)
    e, ds = d + 9, dsig(d)
    while dsig(e) != ds:
        e += 9
    print(e)
