"""Determine the longest common subsequence between two strings."""
import fileinput

def btr(c, x, y, i, j):
    """Backtrack the choices taken when computing the table."""
    if i < 1 or j < 1:
        return ""
    elif x[i-1] == y[j-1]:
        return btr(c, x, y, i-1, j-1) + x[i-1]
    elif c[i][j-1] > c[i-1][j]:
        return btr(c, x, y, i, j-1)
    else:
        return btr(c, x, y, i-1, j)

for test in filter(None, (i.rstrip('\n') for i in fileinput.input())):
    x, y = test.split(';')
    c = [[0 for _ in range(len(y)+1)] for _ in range(len(x)+1)]
    for ix, i in enumerate(x):
        for jx, j in enumerate(y):
            if i == j:
                c[ix+1][jx+1] = c[ix][jx] + 1
            else:
                c[ix+1][jx+1] = max(c[ix+1][jx], c[ix][jx+1])
    print(btr(c, x, y, len(x), len(y)))
