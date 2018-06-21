def normalize(name):
    return name.capitalize()

L1 = ['adam', "LISA", "barI"]
L2 = list(map(normalize, L1))
print(L2)
