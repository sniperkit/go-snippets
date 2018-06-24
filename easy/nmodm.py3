import fileinput

for line in fileinput.input():
    h, l = [int(i) for i in line.split(',')]
    print(h - (h // l) * l)
