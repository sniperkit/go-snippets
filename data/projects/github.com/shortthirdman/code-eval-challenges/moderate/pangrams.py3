import fileinput

for line in fileinput.input(openhook=fileinput.hook_encoded("iso-8859-1")):
    s = line.rstrip('\n').lower()
    r = [i for i in 'abcdefghijklmnopqrstuvwxyz' if i not in s]
    print(''.join(r) if r else "NULL")
