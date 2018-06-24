import fileinput, heapq

for line in fileinput.input():
    n, k, a, b, c, r = [int(i) for i in line.split(',')]
    m = [a]
    for _ in range(k-1):
        m.append((b * m[-1] + c) % r)
    h = [i for i in range(k+1) if i not in m[-k:-1]]
    heapq.heapify(h)
    while len(m)+1 < n:
        o = heapq.heappop(h)
        if not m[-k] in h and not m[-k] in m[-k+1:-1]:
            heapq.heappush(h, m[-k])
        m.append(o)
    print(heapq.heappop(h))
