"""Print out Pascal's triangle up to the requested depth in row major form"""
import fileinput

for line in fileinput.input():
    try:
        n = int(line)
    except ValueError:
        continue
    l = [1]
    st = [1]
    for _ in range(n-1):
        m = [0] + l
        for ix, i in enumerate(l):
            m[ix] += i
        st += m
        l = m
    print(*st)
