"""Determine if the second string is a substring of the first"""
import fileinput

def subcheck(q, r):
    """must not use any substr type library function"""
    while q < len(st[0]) and r < len(st[1]):
        if st[1][r] == '*':
            m = False
            for i in range(q, len(st[0])):
                if subcheck(i, r+1):
                    m = True
                    break
            return m
        elif st[1][r] == '\\' and r < len(st[1])+1 and st[1][r+1] == '*':
            if st[0][q] == '*':
                r += 2
            else:
                return False
        elif st[0][q] != st[1][r]:
            return False
        else:
            r += 1
        q += 1
    return r >= len(st[1])

for line in fileinput.input():
    st = line.split(',')
    st[1] = st[1].lstrip('*').rstrip('\n')
    while len(st[1]) > 1 and st[1][-1] == '*' and st[1][-2] != '\\':
        st[1] = st[1][:-2]

    match = False
    for i in range(len(st[0])):
        if subcheck(i, 0):
            match = True
            break
    print("true" if match else "false")
