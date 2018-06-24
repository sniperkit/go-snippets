for line in io.lines(arg[1]) do
  local a = {}
  for i in line:gmatch("%S+") do
    local x, y, z = i:match("(%d)(%a+)(%d)")
    a[#a+1] = z .. y .. x
  end
  print(table.concat(a, " "))
end
