import fileinput

for line in (i.rstrip('\n') for i in fileinput.input()):
    if len(line) > 55:
        i, p = -1, 40
        while True:
            i = line.find(' ', i+1, 40)
            if i < 0:
                break
            p = i
        line = line[:p] + "... <Read More>"
    print(line)
