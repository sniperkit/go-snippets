classmates=["Michael", "Bob", "Tracy"]
print(classmates)
print(len(classmates))
print(classmates[0])
print(classmates[1])
print(classmates[2])
print(classmates[-1])
print(classmates[-2])
print(classmates[-3])

classmates.append("Adam")
classmates.append("Jordan")
print(classmates)

classmates.insert(1, "Jack")
print(classmates)

print(classmates.pop())
print(classmates)

print(classmates.pop(1))
print(classmates)

classmates[2]="liwei"
print(classmates)

L = ["Apple", 123, True]
print(L)

Ls = ["python", "Java", ["asp", "php"], "scheme"]
print(Ls)
print(Ls[2][0])
print(Ls[2][1])

a = ["f", "e", "a", "g", "d"]
a.sort()
print(a)
