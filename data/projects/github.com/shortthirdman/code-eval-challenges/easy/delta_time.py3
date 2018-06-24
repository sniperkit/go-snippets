import fileinput

for line in fileinput.input():
    s = [int(i) for i in line.replace(":", " ").split()]
    t = abs((s[0] - s[3])*3600 + (s[1] - s[4])*60 + s[2] - s[5])
    print("{:0>2d}".format(t//3600) + ":" + "{:0>2d}".format((t//60)%60) + ":" + "{:0>2d}".format(t%60))
