import fileinput

units = ['PENNY', 'NICKEL', 'DIME', 'QUARTER', 'HALF DOLLAR', 'ONE',
         'TWO', 'FIVE', 'TEN', 'TWENTY', 'FIFTY', 'ONE HUNDRED']
value = [1, 5, 10, 25, 50, 100, 200, 500, 1000, 2000, 5000, 10000]

def centi(a):
    if not '.' in a:
        return int(a) * 100
    a1, a2 = [int(i) for i in a.split('.')]
    return a1 * 100 + a2

for line in fileinput.input():
    p, c = [centi(i) for i in line.split(';')]
    r = []
    if p > c:
        print('ERROR')
        continue
    elif p == c:
        print('ZERO')
        continue
    for i in range(len(units) - 1, -1, -1):
        while p + value[i] <= c:
            r.append(units[i])
            p += value[i]
    print(','.join(r))
