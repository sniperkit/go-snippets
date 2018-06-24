import fileinput

for line in fileinput.input():
    st, r = line.rstrip('\n'), 0
    a = [0] * 26
    for i in st:
        if i.isalpha():
            a[ord(i.lower())-97] += 1
    a.sort()
    for i in range(25, -1, -1):
        if a[i] == 0:
            break
        r += (i+1) * a[i]
    print(r)
