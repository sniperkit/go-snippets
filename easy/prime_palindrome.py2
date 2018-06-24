mp = 0
n = 1000
s = [True] * (n / 2)
for i in xrange(3, int(n ** 0.5) + 1, 2):
    if s[i / 2]:
        s[i * i / 2::i] = [False] * ((n - i * i - 1) / (2 * i) + 1)
for i in xrange(1, n / 2):
    p = 2 * i + 1
    if s[i] and str(p) == str(p)[::-1]:
        mp = p
print mp
