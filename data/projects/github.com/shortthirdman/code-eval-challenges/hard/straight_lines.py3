import fileinput

class Point:
    def __init__(self, x=0, y=0):
        self.x = x
        self.y = y

def in_line(a, b, c):
    return (a.x - b.x) * (a.y - c.y) == (a.x - c.x) * (a.y - b.y)

def is_line(ax, bx, t):
    for ix, i in enumerate(t):
        if ix != ax and ix != bx and in_line(i, t[ax], t[bx]):
            return ix >= bx
    return False

for line in fileinput.input():
    s = [Point(*[int(j) for j in i.split()]) for i in line.split(" | ")]
    print(sum(1 for i in range(len(s)-2) for j in range(i+1, len(s)-1) if is_line(i, j, s)))
