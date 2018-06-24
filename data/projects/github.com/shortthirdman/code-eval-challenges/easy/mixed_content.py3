import fileinput

for line in fileinput.input():
    st, r = line.rstrip('\n').split(','), ""
    n, s = [i for i in st if i.isdigit()], [i for i in st if not i.isdigit()]
    if len(s) > 0:
        r = ','.join(s)
        if len(n) > 0:
            r += "|"
    if len(n) > 0:
        r += ','.join(n)
    print(r)
