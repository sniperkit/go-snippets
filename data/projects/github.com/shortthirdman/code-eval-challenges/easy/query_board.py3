import fileinput

m = [[0] * 256 for _ in range(256)]
for line in fileinput.input():
    s = line.split()
    s1 = int(s[1])
    if len(s) > 2:
        s2 = int(s[2])
    if s[0] == "SetRow":
        m[s1] = [s2] * 256
    elif s[0] == "SetCol":
        for i in range(256):
            m[i][s1] = s2
    elif s[0] == "QueryRow":
        print(sum(m[s1]))
    elif s[0] == "QueryCol":
        print(sum(m[i][s1] for i in range(256)))
