"""predict the next word using the n-gram algorithm"""
import fileinput, string

TEXT = """
Mary had a little lamb its fleece was white as snow;
And everywhere that Mary went, the lamb was sure to go. 
It followed her to school one day, which was against the rule;
It made the children laugh and play, to see a lamb at school.
And so the teacher turned it out, but still it lingered near,
And waited patiently about till Mary did appear.
"Why does the lamb love Mary so?" the eager children cry;
"Why, Mary loves the lamb, you know" the teacher did reply."
"""

TXL = TEXT.split()
for ix, i in enumerate(TXL):
    TXL[ix] = i.strip(string.punctuation)

for line in fileinput.input():
    st = line.split(',')
    n = int(st[0])
    x = st[1].split()

    tnl = []
    for i in range(len(TXL) - n + 1):
        tnl.append(tuple(TXL[i + j] for j in range(n)))

    r = []
    for i in tnl:
        f = True
        for j in range(len(x)):
            if i[j] != x[j]:
                f = False
        if f:
            r.append(i)
    r = sorted(r, key=lambda s: (-r.count(s), ' '.join(s)))
    for ix, i in enumerate(r):
        if ix == 0 or i != r[ix - 1]:
            if ix > 0:
                print(';', end = '')
            print(' '.join(i[len(x):]), end = ',')
            print("%.3f" % (float(r.count(i))/float(len(r))), end = '')
    print()
