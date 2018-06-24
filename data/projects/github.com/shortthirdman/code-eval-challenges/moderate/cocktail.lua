function div(a, b) return math.floor(a/b) end

for line in io.lines(arg[1]) do
  local a, n = {}, 0
  for i in line:gmatch("%d+") do
    a[#a+1] = tonumber(i)
  end
  n, a[#a] = math.min(a[#a], div(#a-1, 2)), nil
  for i = 1, n do
    for j = 2, #a do
      if a[j] < a[j-1] then a[j], a[j-1] = a[j-1], a[j] end
    end
    for j = #a-1, 2, -1 do
      if a[j] < a[j-1] then a[j], a[j-1] = a[j-1], a[j] end
    end
  end
  print(table.concat(a, " "))
end
