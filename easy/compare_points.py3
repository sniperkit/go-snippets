import fileinput

for line in fileinput.input():
    x1, y1, x2, y2 = [int(i) for i in line.split()]
    print((x1 == x2) and ((y1 == y2) and "here" or (y1 < y2) and "N" or "S") or (y1 == y2) and ((x1 < x2) and "E" or "W") or ((x1 < x2) and ((y1 < y2) and "NE" or "SE") or ((y1 < y2) and "NW" or "SW")))
