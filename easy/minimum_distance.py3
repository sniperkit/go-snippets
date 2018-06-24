import fileinput

for line in fileinput.input():
    s = [int(i) for i in line.split()][1:]
    s.sort()
    u = s[len(s)//2]
    print(sum(abs(i - u) for i in s))
