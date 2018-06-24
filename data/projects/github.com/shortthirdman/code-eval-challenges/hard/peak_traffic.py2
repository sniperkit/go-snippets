import collections, fileinput

def BronKerbosch(R, P, X):
    if not P and not X:
        pl = []
        for i in R:
            pl.append(addr[i])
        pl.sort()
        po.append(pl)
    for i in P:
        BronKerbosch(R | set([i]), P & set(g[i]), X & set(g[i]))
        P = P - set([i])
        X = X | set([i])

user = {}
addr = {}
po = []
n = 0
g = collections.defaultdict(list)
q = collections.defaultdict(list)
u = set()

for line in fileinput.input():
    st = line.split('\t')

    a, b = st[1], st[2].strip()
    if a not in user:
        user[a] = n
        addr[n] = a
        n += 1
    if b not in user:
        user[b] = n
        addr[n] = b
        n += 1
    a, b = user[a], user[b]

    if a in g and b in g[a]:
        continue

    if a in q and b in q[a]:
        g[a].append(b)
        g[b].append(a)
        u.add(a)
        u.add(b)
        q[a].remove(b)
        continue
    q[b].append(a)

BronKerbosch(set(), u, set())

for ix, i in enumerate(po):
    po[ix] = ', '.join(j for j in i)
po.sort()
for i in po:
    print i
