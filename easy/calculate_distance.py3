import fileinput, math

for line in fileinput.input():
    s = [int(i) for i in line.translate({ord(i):None for i in '(),'}).split()]
    x, y = s[0]-s[2], s[1]-s[3]
    print(int(math.sqrt(x*x+y*y)))
