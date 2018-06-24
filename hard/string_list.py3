"""Print all of the possible ways to write a string of length N from
the characters in string S, comma delimited in alphabetical order."""
import fileinput, itertools

for line in fileinput.input():
    s = line.rstrip('\n').split(',')
    n = int(s[0])
    b = {''.join(i) for i in itertools.permutations(''.join(sorted(set(s[1])) * n), n)}
    print(*sorted(b), sep=',')
