import fileinput

for line in fileinput.input():
    s = [i.split() for i in line.split(" | ")]
    if sum([int(i, 16) for i in s[0]]) > sum([int(i, 2) for i in s[1]]):
        print("False")
    else:
        print("True")
