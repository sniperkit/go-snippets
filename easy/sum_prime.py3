def isp(a, p):
    for j in p:
        if a % j == 0:
            return False
    return True

n = 1000
p = [2]
i = 3
sp = 2
while len(p) < n:
    if isp(i, p):
        p += [i]
        sp += i
    i += 2
print(sp)
