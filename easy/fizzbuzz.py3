import fileinput

for line in fileinput.input():
    f, b, n = [int(i) for i in line.split()]
    print(" ".join(["FB" if i % f == 0 and i % b == 0 else "F" if i % f == 0 else "B" if i % b == 0 else str(i) for i in range(1, n+1)]))
