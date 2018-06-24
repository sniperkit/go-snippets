for line in io.lines(arg[1]) do
  local a, sep = {}, line:find("|")
  local m = tonumber(line:sub(sep+1):match("%d+"))
  for i in line:sub(1, sep-1):gmatch("%S+") do
    a[#a+1] = i
  end
  while #a > 1 do
    local n = m % #a
    table.remove(a, (n > 0) and n or #a)
  end
  print(a[1])
end
