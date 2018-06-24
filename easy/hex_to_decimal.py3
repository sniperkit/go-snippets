import fileinput

print('\n'.join([str(int(i, 16)) for i in fileinput.input()]))
