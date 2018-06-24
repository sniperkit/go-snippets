import fileinput

for line in fileinput.input():
    a, num, ar, rr = 0, 0, 0, 0
    for x in line.rstrip('\n'):
        if x.isdigit():
            a = int(x)
        else:
            if x == 'M':
                r = 1000
            elif x == 'D':
                r = 500
            elif x == 'C':
                r = 100
            elif x == 'L':
                r = 50
            elif x == 'X':
                r = 10
            elif x == 'V':
                r = 5
            elif x == 'I':
                r = 1
            if r > rr:
                num -= ar * rr
            else:
                num += ar * rr
            ar, rr = a, r
    print(num + ar * rr)
