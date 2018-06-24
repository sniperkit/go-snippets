import fileinput

for line in fileinput.input():
    print(sum([int(i) for i in line.rstrip('\n')]))
