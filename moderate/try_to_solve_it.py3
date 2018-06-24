import fileinput

s, d = "abcdefghijklmnopqrstuvwxyz", "uvwxyznopqrstghijklmabcdef"
m = {ord(i): ord(d[ix]) for ix, i in enumerate(s)}
for line in fileinput.input():
    print(line.rstrip('\n').translate(m))
