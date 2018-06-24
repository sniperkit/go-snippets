import fileinput

def stairs(a):
    b = c = 1
    for _ in range(a-1):
        b, c = c, b + c
    return c

for line in fileinput.input():
    if not line == '\n':
        print(stairs(int(line)))
