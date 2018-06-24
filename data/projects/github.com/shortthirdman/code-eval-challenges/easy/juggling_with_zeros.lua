for line in io.lines(arg[1]) do
  local r = 0
  for i in line:gmatch("0+ 0+") do
    j = string.find(i, " ")
    r = r * 2^(#i-j)
    if j == 3 then r = r + 2^(#i-j)-1 end
  end
  print(r)
end
