import fileinput

for line in fileinput.input():
    print(' '.join([i[0].upper() + i[1:] for i in line.split()]))
