import fileinput, math

buckets = 7759

sq = [[] for i in range(buckets)]
for i in range(46341):
    i2 = i*i
    sq[i % buckets].append(i2)

for line in fileinput.input():
    if fileinput.isfirstline():
        continue
    x = int(line)
    top = int(math.sqrt(x)) + 1
    bot = int(math.sqrt(x/2))
    l = []
    n = 0
    for i in range(bot, top):
        t = (x - i*i)
        if t in sq[int(math.sqrt(t)) % buckets] and t not in l:
            l.append(i*i)
            n += 1
    print(n)
