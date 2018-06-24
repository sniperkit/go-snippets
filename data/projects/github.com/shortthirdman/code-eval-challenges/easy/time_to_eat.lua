for line in io.lines(arg[1]) do
  local a = {}
  for i in line:gmatch("%S+") do a[#a + 1] = i end
  table.sort(a, function (a,b) return (a > b) end)
  print(table.concat(a, " "))
end
