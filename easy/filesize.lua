local f = assert(io.open(arg[1], "r"))
print(f:seek("end"))
f:close()
