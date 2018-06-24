import fileinput

for line in fileinput.input():
    for c in line:
        if c.islower():
            print(c.upper(), end = '')
        elif c.isupper():
            print(c.lower(), end = '')
        else:
            print(c, end = '')
    if line[-1] != '\n':
        print()
