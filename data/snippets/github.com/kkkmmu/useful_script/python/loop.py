names = ["Michael", "Bob", "Tracy"]
for name in names:
    print(name)

sum = 0
for x in [1, 2, 3, 4, 5, 6, 7, 8, 9]:
    sum += x
print(sum)
print(list(range(20)))

sum = 0
for x in range(100):
    sum += x
print(sum)

sum = 0
n = 99
while n > 0:
    sum += n
    n = n-2
print(sum)

s = range(10)

for i in s:
    print(i)

for ch in "ABCDefg":
    print(ch)

from collections import Iterable

print(isinstance("abc", Iterable))
print(isinstance([1, 2, 3], Iterable))
print(isinstance(123, Iterable))
print(isinstance((1, 2, 3), Iterable))
print(isinstance({"a": 1, "b": 2, "c":3}, Iterable))
print(isinstance((1, 2, 3), Iterable))

print([x*x for x in range (1, 11)])
print([x*x for x in range (1, 11) if x % 2 == 0])
print([m + n for m in "ABC" for n in "abc"])

import os
print([d for d in os.listdir(".")])

g = (x *x for x in range (10))
print(next(g))
print(next(g))
print(next(g))
for n in g:
     print(n)
