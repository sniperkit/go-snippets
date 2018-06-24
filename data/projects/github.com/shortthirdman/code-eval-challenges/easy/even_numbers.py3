import fileinput

print(*[1 - int(i) & 1 for i in fileinput.input()], sep='\n')
