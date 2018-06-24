import fileinput

for line in fileinput.input():
    s = line.split(',')
    print(s[0].translate({ord(i):None for i in s[1].strip()}))
