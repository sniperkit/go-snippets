"""Print out integers in English textual representation"""
import fileinput

STR0 = ['One', 'Two', 'Three', 'Four', 'Five', 'Six', 'Seven', 'Eight', 'Nine']
STR1 = ['Ten', 'Eleven', 'Twelve', 'Thirteen', 'Fourteen', 'Fifteen',
        'Sixteen', 'Seventeen', 'Eighteen', 'Nineteen']
STR2 = ['Twenty', 'Thirty', 'Forty', 'Fifty',
        'Sixty', 'Seventy', 'Eighty', 'Ninety']
STR3 = 'Hundred'
STR4 = 'Thousand'
STR5 = 'Million'

def wrd(num):
    """Convert three digits to text."""
    txt = ''
    if num[0] != '0':
        txt = STR0[int(num[0]) - 1] + STR3
    if num[1] == '0':
        if num[2] != '0':
            txt = txt + STR0[int(num[2]) - 1]
    elif num[1] == '1':
        txt = txt + STR1[int(num[2])]
    else:
        txt = txt + STR2[int(num[1]) - 2]
        if num[2] != '0':
            txt = txt + STR0[int(num[2]) - 1]
    return txt

for line in fileinput.input():
    line = ''.join(c for c in line if c in '0123456789')
    if not line:
        continue
    line = '0' * (9 - len(line)) + line

    s = ''
    t = wrd(line[0:3])
    if t:
        s = s + t + STR5
    t = wrd(line[3:6])
    if t:
        s = s + t + STR4
    t = wrd(line[6:9])
    s = s + t
    print s + 'Dollars'
