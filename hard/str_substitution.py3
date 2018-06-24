"""replace a list of strings in a string"""
import fileinput

for line in fileinput.input():
    st = line.split(';')
    a = [st[0]]
    ls = st[1].split(',')
    ls[-1] = ls[-1].strip()

    lsc = -1
    while ls:
        j = []
        lsc += 2
        for i in a:
            j.append(i)
            if type(i) == str:
                while i.count(ls[0]):
                    del j[-1]
                    j.extend(i.partition(ls[0]))
                    i = j[-1]
                    j[-2] = lsc
        a = j
        del ls[0:2]

    ls = st[1].split(',')
    ls[-1] = ls[-1].strip()

    for ix, i in enumerate(a):
        if type(i) == int:
            a[ix] = ls[i]
    print(*a, sep='')
