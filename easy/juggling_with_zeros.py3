import fileinput, math

for line in fileinput.input():
    st, n = line.split(), 0
    for i in range(len(st)//2):
        t = math.pow(2, len(st[i*2+1]))
        n *= t
        if st[i*2] == "00":
            n += t - 1
    print(int(n))
