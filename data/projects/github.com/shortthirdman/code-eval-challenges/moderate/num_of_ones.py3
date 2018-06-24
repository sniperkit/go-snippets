import fileinput

print(*[bin(int(i)).count("1") for i in fileinput.input()], sep='\n')
