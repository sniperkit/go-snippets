"""Print all of the possible ways to write a string of length N from
the characters in string S, comma delimited in alphabetical order."""
import fileinput, itertools

for line in fileinput.input():
    st = line.split(',')
    n = int(st[0])
    a = st[1].strip()
    b = []
    for s in itertools.permutations(''.join(sorted(set(a)) * n), n):
        c = ''.join(s)
        if c not in b:
            b.append(c)
    print ','.join(sorted(b))
