for line in io.lines(arg[1]) do
  local a, r = {}, {}
  for i in line:gmatch("[^ ;]+") do a[#a+1] = i end
  local h = math.ceil(#a/2)
  for i = 1, h-1 do
    r[tonumber(a[h+i])] = a[i]
  end
  for i = 1, h do
    if not r[i] then r[i] = a[h] break end
  end
  print(table.concat(r, " "))
end
