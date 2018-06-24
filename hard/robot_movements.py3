def robot(f, x, y):
    if x == 3 and y == 3:
        return 1
    ret = 0
    if x > 0 and f & (1 << (x - 1 + 4 * y)) == 0:
        ret += robot(f | (1 << (x - 1 + 4 * y)), x - 1, y)
    if y > 0 and f & (1 << (x + 4 * (y - 1))) == 0:
        ret += robot(f | (1 << (x + 4 * (y - 1))), x, y - 1)
    if x < 3 and f & (1 << (x + 1 + 4 * y)) == 0:
        ret += robot(f | (1 << (x + 1 + 4 * y)), x + 1, y)
    if y < 3 and f & (1 << (x + 4 * (y + 1))) == 0:
        ret += robot(f | (1 << (x + 4 * (y + 1))), x, y + 1)
    return ret

print(robot(1, 0, 0))
