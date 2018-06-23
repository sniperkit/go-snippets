print("you", "are", "a", "good", "man")

print("100+200=", 100+200)

print("\\\t\\")
print(r"\\\t\\")

print('''line1
        line2
        line3''')

print(True and False)
print(True and True)
print(True or False)
print(False or False)

print(None)

a = 123
print(a)
a = "ABC"
print(a)

print(10/3)
print(10//3)

print("包含中文的str")

print(ord("A"))
print(ord("中"))
print(chr(66))
print(chr(88))
print(chr(25991))

print('\u4e2d\u6587')

print("ABC")
print(b"ABC")

print("ABC".encode("ascii"))
print("ABC".encode("utf-8"))
#print("中文".encode("ascii"))
print("中文".encode("utf-8"))

print(b"ABC".decode("ascii"))
print(b"ABC".decode("utf-8"))

print(len("ABC"))
print(len("中文"))
print(len("中文".encode("utf-8")))

print("Hello, %s" % "world")
print("Hi %s, your have $%d" % ("Michael", 10000))
