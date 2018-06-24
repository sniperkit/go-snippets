import fileinput

for line in fileinput.input():
    s = [int(i) for i in line if i.isdigit()]
    for i in range(len(s)-2, -1, -2):
        s[i] *= 2
    print("Real" if sum(s) % 10 == 0 else "Fake")
