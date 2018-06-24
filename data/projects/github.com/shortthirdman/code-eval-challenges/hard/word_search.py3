"""Find if the word exists in the grid."""
import fileinput

BOARD = ["ABCE", "SFCS", "ADEE"]

def fpath(pos, bee, st0):
    """Find a path."""
    if not st0:
        return True
    neigh = []
    if pos[0] > 0 and bee[pos[0]-1][pos[1]] == st0[0]:
        neigh.append([pos[0]-1, pos[1]])
    if pos[0] < 2 and bee[pos[0]+1][pos[1]] == st0[0]:
        neigh.append([pos[0]+1, pos[1]])
    if pos[1] > 0 and bee[pos[0]][pos[1]-1] == st0[0]:
        neigh.append([pos[0], pos[1]-1])
    if pos[1] < 3 and bee[pos[0]][pos[1]+1] == st0[0]:
        neigh.append([pos[0], pos[1]+1])
    if not neigh:
        return False
    for nei in neigh:
        be0 = list(bee)
        be0[pos[0]] = bee[pos[0]][:pos[1]] + " " + bee[pos[0]][pos[1]+1:]
        if fpath(nei, be0, st0[1:]):
            return True
    return False

for line in filter(None, (i.rstrip('\n') for i in fileinput.input())):
    start = []
    for ix, i in enumerate(BOARD):
        for jx, j in enumerate(i):
            if line[0] == j:
                start.append([ix, jx])
    f = False
    for i in start:
        if fpath(i, list(BOARD), line[1:]):
            f = True
            break
    print("True" if f else "False")
