import fileinput

k = r""" !"#$%&'()*+,-./0123456789:<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"""
for line in fileinput.input():
    s = line.rstrip('\n').split(';')
    for ix, i in enumerate(s[1]):
        print(k[(len(k) + k.find(i) - int(s[0][ix % len(s[0])])) % len(k)], end='')
    print()
