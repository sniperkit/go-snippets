import fileinput

p = q = -1
for line in fileinput.input():
    s = list(line.rstrip('\n'))
    if p < 0:
        try:
            p = s.index('C')
        except ValueError:
            p = s.index('_')
        s[p] = "|"
    else:
        while True:
            try:
                q = s.index('C')
            except ValueError:
                q = -1
                break
            if abs(p-q) < 2:
                break
        if q < 0:
            while True:
                q = s.index('_')
                if abs(p-q) < 2:
                    break
        if q < p:
            sym = "/"
        elif q == p:
            sym = "|"
        else:
            sym = "\\"
        s[q] = sym
        p = q
    print("".join(s))
