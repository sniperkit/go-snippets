import fileinput

def isBalanced(s, c):
    while c >= 0:
        while s != "" and s[0] not in "():":
            s = s[1:]
        while s != "" and s[-1] not in "()":
            s = s[:-1]
        if s == "":
            return c == 0
        if s[0] == '(':
            s, c = s[1:], c+1
        elif s[0] == ')':
            s, c = s[1:], c-1
        elif s[0] == ':':
            if s[1] == '(':
                return isBalanced(s[2:], c) or isBalanced(s[2:], c+1)
            elif s[1] == ')':
                return isBalanced(s[2:], c) or isBalanced(s[2:], c-1)
            else:
                s = s[1:]
        else:
            return False
    return False

for line in fileinput.input():
    print("YES" if isBalanced(line, 0) else "NO")
