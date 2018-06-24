import fileinput

for line in filter(lambda x: x != '\n', fileinput.input()):
    s = line.split("| ")
    print(*[s[0][int(i)-1] for i in s[1].split()], sep='')
