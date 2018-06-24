import fileinput

def isPalindrome(a):
    return str(a) == str(a)[::-1]

for line in fileinput.input():
    l, r = [int(i) for i in line.split()]
    n = 0
    for i in range(l, r+1):
        prev = -1
        for j in range(i, r+1):
            if prev >= 0:
                p = prev
                if isPalindrome(j):
                    p += 1
            else:
                p = 0
                for k in range(i, j+1):
                    if isPalindrome(k):
                        p += 1
            if not p & 1:
                n += 1
            prev = p
    print(n)
