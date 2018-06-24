import fileinput

def g2d(a):
    a ^= a >> 4
    a ^= a >> 2
    a ^= a >> 1
    return a

for line in fileinput.input():
    print(" | ".join([str(g2d(int(i, 2))) for i in line.split(" | ")]))
