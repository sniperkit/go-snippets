import fileinput

def subs(s, t):
    if len(t) == 0:
        return 1
    if len(s) == 0:
        return 0
    if s[0] == t[0]:
        return subs(s[1:], t[1:]) + subs(s[1:], t)
    return subs(s[1:], t)

for line in fileinput.input():
    st = line.rstrip('\n').split(",")
    print(subs(st[0], st[1]))
