import fileinput

for line in fileinput.input():
    s = [ord(i) for i in line.rstrip('\n')]
    k = 0
    for ix, i in enumerate(s):
        if i == ord('$'):
            k = ix
        s[ix] = (s[ix] << 16) + ix
    s.sort()
    for _ in range(len(s)-1):
        print(chr(s[k] >> 16), end='')
        k = s[k] % (1 << 16)
    print()
