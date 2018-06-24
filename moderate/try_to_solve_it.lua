function deco(c, o)
  return string.char(string.byte(c) + o)
end

for line in io.lines(arg[1]) do
  for i in line:gmatch(".") do
    io.write(i:find("[a-f]") and deco(i, 20) or
             i:find("[g-m]") and deco(i, 7) or
             i:find("[n-t]") and deco(i, -7) or
             i:find("[u-z]") and deco(i, -20) or i)
  end
  print()
end
