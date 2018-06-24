import fileinput

for line in fileinput.input():
    s = [int(i) for i in line.split()]
    l, d, n, count, t, tx, i = s[0], s[1], s[2], 0, 0, 6-s[1], 6
    while i <= l-6:
        if i > tx-d:
            i = tx
            if t == n:
                tx = l-6+d
            else:
                tx, t = s[t+3], t + 1
        else:
            count = count + 1
        i += d
    print(count)
