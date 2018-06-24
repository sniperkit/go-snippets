import fileinput

for line in fileinput.input():
    s = [float(i) for i in line.translate({ord(i):None for i in 'Centr:RadiusPo(),;'}).split()]
    x, y, r = s[3] - s[0], s[4] - s[1], s[2]
    print("true" if x*x + y*y <= r*r else "false")
