s = set([1, 1, 1, 2, 3, 3, 3, 4, 5])
print(s)
s.add(10)
s.add(8)
s.add(8)
s.add(8)
s.add(8)
print(s)

s.remove(8)
print(s)

s1 = set([1, 2, 3, 4])
s2 = set([4, 5, 6, 7])

print(s1 & s2)
print(s1 | s2)
