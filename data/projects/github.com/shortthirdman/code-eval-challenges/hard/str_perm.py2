"""print out all the permutations of a string in alphabetical order"""
import fileinput, itertools

for line in fileinput.input():
    st = ''
    for s in itertools.permutations(sorted(line.strip())):
        st = st + ''.join(s) + ','
    print st.rstrip(',')
