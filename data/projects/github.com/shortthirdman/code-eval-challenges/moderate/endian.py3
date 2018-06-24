import sys

print("BigEndian" if sys.byteorder == "big" else "LittleEndian")
