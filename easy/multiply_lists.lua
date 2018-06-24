for line in io.lines(arg[1]) do
  local a = {}
  for i in line:gmatch("%d+") do
    a[#a + 1] = i
  end
  for i = 1, #a/2 do
    if i > 1 then io.write " " end
    io.write(a[i] * a[#a/2 + i])
  end
  print()
end
