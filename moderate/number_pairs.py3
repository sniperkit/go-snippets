import fileinput

for line in fileinput.input():
    st = line.split(';')
    a = st[0].split(',')
    b = int(st[1])

    f = False
    for i in a:
        if 2 * int(i) < b and str(b - int(i)) in a:
            if f:
                print(';', end='')
            print(i, ',', b - int(i), sep='', end='')
            f = True
    if not f:
        print('NULL', end='')
    print()
