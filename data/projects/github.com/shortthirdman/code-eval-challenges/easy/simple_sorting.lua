for line in io.lines(arg[1]) do
  a = {}
  for i in line:gmatch("%S+") do a[#a + 1] = tonumber(i) end
  table.sort(a)
  for ix, i in ipairs(a) do
    if ix > 1 then io.write(" ") end
    io.write(string.format("%.3f", i))
  end
  print()
end
