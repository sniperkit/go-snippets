import fileinput, sys

morse = "ETIANMSURWDKGOHVF L PJBXCYZQ  54 3   2       16       7   8 90"
for line in fileinput.input():
    m = 1
    for i in line.rstrip('\n'):
        if i == '.':
            m *= 2
        elif i == '-':
            m = m*2 + 1
        else:
            if m == 1:
                print(' ', end='')
            else:
                if m < 64:
                    print(morse[m-2], end='')
                m = 1
    print(morse[m-2] if m > 1 and m < 64 else '')
