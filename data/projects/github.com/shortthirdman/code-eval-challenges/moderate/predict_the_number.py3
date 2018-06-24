import fileinput

print(*[bin(int(i)).count("1") % 3 for i in fileinput.input()], sep='\n')
