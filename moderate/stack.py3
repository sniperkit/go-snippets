import fileinput

def stackpush(a):
    stack.append(a)

def stackpop():
    return stack.pop()

for line in fileinput.input():
    stack = []
    st = line.split()

    for i in st:
        stackpush(int(i))

    if stack:
        print(stackpop(), end='')
        if stack:
            stackpop()
    while stack:
        print(' ' + str(stackpop()), end='')
        if stack:
            stackpop()
    print()
