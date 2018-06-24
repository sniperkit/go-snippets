import fileinput

print(*["{0:b}".format(int(i)) for i in fileinput.input() if i != '\n'], sep='\n')
