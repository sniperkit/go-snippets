import fileinput

for line in fileinput.input():
    st = line.split()
    st.sort(key=float)
    print(*st)
