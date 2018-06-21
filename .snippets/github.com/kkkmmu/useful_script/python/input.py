name = input()
print("Hello", name)

name = input("Please input your name: ")
print("Hello", name)

#birth = input("birth: ")

#if birth < 2000:
#    print("good")
#else:
#    print("sb")

s = input("birth: ")
birth = int(s)
if birth < 2000:
    print("good")
else:
    print("sb")

#BMI

hs = input("Please input your height: ")
h = int(hs)
ws = input("Please input your weight: ")
w = float(ws)
bmi = w / (h*h)

if bmi < 18.5:
    print("Your weight is too low")
elif bmi < 25 and bmi > 18.5:
    print("Your weight is normal")
elif bmi >= 25 and bmi <= 32:
    print("You are a little fat")
else:
    print("You are too fat")

