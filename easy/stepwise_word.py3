import fileinput, functools

for line in fileinput.input():
    l = functools.reduce(lambda a, b: b if len(b) > len(a) else a, line.split())
    print(*[ix * '*' + i for ix, i in enumerate(l)])
