import fileinput

for line in fileinput.input():
    n = int(line)
    if n > 2:
        print(2, end='')
    s = [True] * int(n / 2)
    for i in range(3, int(n ** 0.5) + 1, 2):
        if s[int(i / 2)]:
            s[int(i * i / 2)::i] = [False] * int((n - i * i - 1) / (2 * i) + 1)
    for i in range(1, int(n / 2)):
        if s[i]:
            print(',', 2 * i + 1, sep='', end='')
    print()
