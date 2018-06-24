import fileinput

for line in filter(None, [i.rstrip('\n') for i in fileinput.input()]):
    s = line.split(',')
    print(1 if s[0].endswith(s[1]) else 0)
