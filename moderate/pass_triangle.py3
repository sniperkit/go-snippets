import fileinput

a = []
for line in fileinput.input():
    s = [int(i) for i in line.split()]
    if len(a) > 0:
        s[0] += a[0]
        if len(a) > 1:
            s[-1] += a[-1]
    for i in range(1, len(s)-1):
        s[i] += max(a[i-1], a[i])
    a = s
print(max(a))
