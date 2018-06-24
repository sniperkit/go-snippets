import fileinput

print(*[i.split()[-2] for i in fileinput.input()], sep='\n')
