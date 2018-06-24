import fileinput

for line in fileinput.input():
    a = float(line)
    b = int(a)
    print(b, '.', sep = '', end = '')
    a = (a - b) * 60
    b = int(a)
    print("{:0>2d}'".format(b), end = '')
    a = (a - b) * 60
    b = int(a)
    print("{:0>2d}\"".format(b))
