for line in io.lines(arg[1]) do
  local a, r = {}, {}
  for i in line:gmatch("-?%d+") do
    a[#a + 1] = tonumber(i)
  end
  local n = select(2, line:gsub("|", "|")) + 1
  local m = #a/n
  for i = 1, m do
    r[i] = a[i]
  end
  for i = 2, n do
    for j = 1, m do
      if a[(i-1)*m + j] > r[j] then
        r[j] = a[(i-1)*m + j]
      end
    end
  end
  print(table.concat(r, " "))
end
