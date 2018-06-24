import fileinput, functools

for line in fileinput.input():
    s = line.rstrip('\n')
    print(functools.reduce(lambda x, y: [x[1]+1, x[1]+1] if y != s[x[1]-x[0]] else [x[0], x[1]+1], s, [1, 0])[0])
