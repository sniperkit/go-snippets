import fileinput

ronum = [1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1]
rostr = ["M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"]
for line in fileinput.input():
    i, num = 0, int(line)
    while num > 0:
        if num >= ronum[i]:
            print(rostr[i], end = '')
            num -= ronum[i]
        else:
            i += 1
    print()
