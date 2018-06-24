import fileinput

def within_rect(rect, pt):
    return pt[0] >= rect[0] and pt[0] <= rect[2] and pt[1] >= rect[3] and pt[1] <= rect[1]

for line in fileinput.input():
    s = [int(i) for i in line.split(',')]
    r1, r2 = [s[0], s[1], s[2], s[3]], [s[4], s[5], s[6], s[7]]
    p1 = [[r1[0], r1[1]], [r1[0], r1[3]], [r1[2], r1[1]], [r1[2], r1[3]]]
    p2 = [[r2[0], r2[1]], [r2[0], r2[3]], [r2[2], r2[1]], [r2[2], r2[3]]]

    f = False
    for i in p1:
        if within_rect(r2, i):
            f = True
            break
    if not f:
        for i in p2:
            if within_rect(r1, i):
                f = True
                break
    print(f)
