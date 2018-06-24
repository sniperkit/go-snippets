import fileinput

def contains(w, x):
    while len(w) > 0 and len(x) > 0:
        if w[0] < x[0]:
            w = w[1:]
        elif w[0] == x[0]:
            w, x = w[1:], x[1:]
        else:
            return False
    return len(x) == 0

for line in fileinput.input():
    st = line.rstrip('\n').split(" | ")
    s, t = st[0].split(), sorted(st[1])
    r = [i for i in s if contains(sorted(i), t)]
    print("False" if len(r) == 0 else ' '.join(r))
