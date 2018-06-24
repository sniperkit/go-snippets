import fileinput, functools

for line in fileinput.input():
    print(functools.reduce(lambda a, b: b if len(b) > len(a) else a, line.rstrip('\n').split()))
