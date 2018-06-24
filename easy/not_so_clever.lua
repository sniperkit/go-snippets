for line in io.lines(arg[1]) do
  local a, n, p, l = {}, 0, 0, 2
  for i in line:gmatch("%d+") do
    a[#a+1] = tonumber(i)
  end
  n, a[#a] = a[#a], nil
  for i = 1, n do
    p, l = l, 0
    for j = p, #a do
      if a[j] < a[j-1] then
        a[j], a[j-1] = a[j-1], a[j]
        l = j > 2 and j - 1 or 3
        break
      end
    end
    if l == 0 then break end
  end
  print(table.concat(a, " "))
end
