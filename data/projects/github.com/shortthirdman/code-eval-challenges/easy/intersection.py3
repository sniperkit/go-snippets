import fileinput

for line in fileinput.input():
    s = [{int(j) for j in i.split(',')} for i in line.rstrip('\n').split(';')]
    print(*sorted(list(s[0] & s[1])), sep=',')
