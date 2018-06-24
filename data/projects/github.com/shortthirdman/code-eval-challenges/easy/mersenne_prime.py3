import fileinput

def isp(a, p):
    for j in p:
        if a % j == 0:
            return False
    return True

p = [2]
i = 3

for line in fileinput.input():
    n = int(line)
    while 2**i-1 < n:
        if isp(i, p):
            p += [i]
        i += 2
    print(", ".join([str(2**j-1) for j in p if 2**j-1 < n]))
