import fileinput, functools

for line in fileinput.input():
    s = line.rstrip('\n')
    print("True" if int(s) == functools.reduce(lambda x, y: x + pow(int(y), len(s)), s, 0) else "False")
