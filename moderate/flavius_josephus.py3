import fileinput

for line in fileinput.input():
    n, m = [int(i) for i in line.split(",")]
    x, l, r = 0, list(range(n)), []
    while len(l) > 0:
        x = (x - 1 + m) % len(l)
        r.append(str(l.pop(x)))
    print(' '.join(r))
