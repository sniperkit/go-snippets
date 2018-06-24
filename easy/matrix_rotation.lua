for line in io.lines(arg[1]) do
  local a, r = {}, {}
  for i in line:gmatch("%S+") do
    a[#a+1] = i
  end
  local l = math.sqrt(#a)
  for i = 1, l do
    for j = l, 1, -1 do
      r[#r+1] = a[(j-1)*l + i]
    end
  end
  print(table.concat(r, " "))
end
