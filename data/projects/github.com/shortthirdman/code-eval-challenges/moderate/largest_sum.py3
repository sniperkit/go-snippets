import fileinput, functools

for line in fileinput.input():
    s = [int(i) for i in line.split(",")]
    print(functools.reduce(lambda a, b: (max(a[0], b+a[1]), max(0, b+a[1])), s, (s[0], 0))[0])
