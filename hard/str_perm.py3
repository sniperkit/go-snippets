"""print out all the permutations of a string in alphabetical order"""
import fileinput, itertools

for line in fileinput.input():
    print(*sorted({''.join(i) for i in itertools.permutations(line.rstrip('\n'))}), sep=',')
