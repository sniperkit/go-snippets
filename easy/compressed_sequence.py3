import fileinput

for line in fileinput.input():
    prev, count, r = "", 0, []
    for i in line.split():
        if i == prev:
            count += 1
        else:
            if count > 0:
                r.extend([str(count), prev])
            count, prev = 1, i
    r.extend([str(count), prev])
    print(*r)
