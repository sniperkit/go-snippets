import fileinput

print(*[i.rstrip('\n').lower() for i in fileinput.input()], sep='\n')
