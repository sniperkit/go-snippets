import fileinput

for line in fileinput.input():
    st = line.split()
    if st.pop(0) == '1':
        print("Jolly")
        continue
    s = sorted(abs(int(st[i]) - int(st[i+1])) for i in range(len(st) - 1))
    print("Jolly" if s == list(range(1, len(st))) else "Not jolly")
