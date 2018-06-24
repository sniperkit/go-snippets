import fileinput, math

for line in fileinput.input():
    a, b = [int(i) for i in line.split(",")]
    c = a - int(math.floor(a * math.pow(2, -math.log(b, 2))) * math.pow(2, math.log(b, 2)))
    print(a - c + b if c > 0 else a)
