"""decode a message that was encoded and sent in two parts"""
import fileinput, itertools

def genkey():
    """generate sequence 0,00,01,10,000,001,010,011,100,101,110,0000,..."""
    sll = 1
    while True:
        for iii in list(itertools.product('01', repeat=sll))[:-1]:
            ret = ""
            for jjj in iii:
                ret += jjj
            yield ret
        sll += 1


MESS = ""
TRANS = {}
XX = genkey()
ST = ""

for line in fileinput.input():
    ST += line.rstrip('\n')

while ST:
    IX, i = 0, ""
    for IX, i in enumerate(ST):
        if i == '0' or i == '1':
            break
        sym = next(XX)
        TRANS[sym] = i
    SEGW = eval('0b' + i + ST[IX+1] + ST[IX+2])
    IX += 3
    while SEGW > 0:
        CST = ""
        for _ in range(SEGW):
            CST += ST[IX]
            IX += 1
        if CST == '1' * SEGW:
            SEGW = eval('0b' + ST[IX] + ST[IX+1] + ST[IX+2])
            IX += 3
            continue
        MESS += TRANS[CST]
    print MESS

    ST = ST[IX:]
    MESS = ""
    TRANS = {}
    XX = genkey()
