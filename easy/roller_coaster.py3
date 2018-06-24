import fileinput

for line in fileinput.input():
    up = True
    for c in line:
        if c.isalpha():
            if up:
                print(c.upper(), end = '')
                up = False
            else:
                print(c.lower(), end = '')
                up = True
        else:
            print(c, end = '')
    if line[-1] != '\n':
        print()
