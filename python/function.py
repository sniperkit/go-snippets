print(abs(-110))
print(max(1, 5, -10, 10))
print(int("123"))
print(int(1.24))
print(float("1.2345"))

print(str(1.23))
print(str(1000))
print(bool(100))
print(bool(-1000))
print(bool(0))
print(bool(None))

a = abs

print(a(-1234))


print(hex(1234))

def mabs(x):
    if x >= 0:
        return x
    else:
        return -x

print(mabs(-11111))

def nop():
    pass

print(nop())

def mabs2(x):
    if not isinstance(x, (int, float)):
        raise TypeError("Bad operand type")
    if x >= 0:
        return x
    else:
        return -x

print(mabs2(-12343))

import math

def move(x, y, step, angle=0):
    nx = x + step * math.cos(angle)
    ny = y - step * math.sin(angle)
    return nx, ny

print(move(10, 10, 1, 30))

z = move(16, 16, 2, 45)
print(z)

def power(x, n=2):
    s = 1
    while n > 0:
        n = n-1
        s = s*x
    return s

print(power(2))
print(power(2, 4))

def defaultp(a, b, c=1, d=2):
    print(a, b, c, d)

defaultp(1, 2)
defaultp(1, 2, d=10)

def calc(*numbers):
    sum = 0
    for n in numbers:
        sum += n*n
    return sum

print(calc(1,2,3,4,5))
print(calc(1,2))
print(calc(1))
print(calc())

nums = [4, 5, 5, 6, 8]
print(calc(*nums))

def person(name, age, **kv):
    print("name:", name, "age:", age, "other:", kv)

person("Liwei", 30)
person("Liwei", 30, city="Beijing")
person("Liwei", 30, city="Beijing", gender="M")

extra = {"City": "Beijing", "job": "polan"}
person("Jack", 100, **extra)

def person2(name, age, *, city, job):
    print(name, age, city, job)

#person2("jack", 24)
#person2("jack", 24, "Beijing", "polan")
person2("jack", 24, city="Beijin", job="polan")

def fact(n):
    if n == 1:
        return 1
    return n * fact(n -1)

print(fact(10))

def f (x):
    return x * x

r = map(f, [1,2,3,4,5,6,7,8,9,10])
print(r)
print(list(r))
print(list(map(str, [1,2,3,4,5,6,7,8,9])))
