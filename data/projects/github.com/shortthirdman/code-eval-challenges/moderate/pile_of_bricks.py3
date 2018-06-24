import fileinput, re

for line in fileinput.input():
    r = []
    s = [int(i) for i in re.findall(r"-?\d+", line)]
    h = sorted([abs(s[0] - s[2]), abs(s[1] - s[3])])
    for i in range(4, len(s), 7):
        b = sorted([abs(s[i+1] - s[i+4]), abs(s[i+2] - s[i+5]), abs(s[i+3] - s[i+6])])
        if b[0] <= h[0] and b[1] <= h[1]: # or diagonal(b, h)
            r.append(s[i])
    print("-" if len(r) == 0 else ','.join([str(i) for i in sorted(r)]))
