import fileinput

phrases = [", yeah!",
           ", this is crazy, I tell ya.",
           ", can U believe this?",
           ", eh?",
           ", aw yea.",
           ", yo.",
           "? No way!",
           ". Awesome!"]
c, l = 0, False
for line in fileinput.input():
    for i in line.rstrip('\n'):
        if i in ".!?":
            if l:
                print(phrases[c], end='')
                c, l = (c+1) % 8, False
            else:
                print(i, end='')
                l = True
        else:
            print(i, end='')
    print()
