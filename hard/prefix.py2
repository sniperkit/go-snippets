"""Evaluate prefix expressions."""
import fileinput

for line in fileinput.input():
    st = line.split(' ')
    n = len(st)
    st[n/2] = ''.join(j for j in st[n/2] if j.isdigit())
    e = st[n/2]
    for i in range(n/2 + 1, len(st)):
        st[n - i - 1] = ''.join(j for j in st[n - i - 1] if j in '+*/')
        st[i] = ''.join(j for j in st[i] if j.isdigit())
        e += st[n - i - 1] + st[i]
        e = str(eval(e))
    print e
