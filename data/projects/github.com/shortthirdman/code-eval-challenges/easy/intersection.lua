for line in io.lines(arg[1]) do
  local r, s = line:match("([^;]+);([^;]+)")
  local a, b, first = {}, {}, true
  for i in r:gmatch("[-]?%d+") do
    a[#a + 1] = tonumber(i)
  end
  for i in s:gmatch("[-]?%d+") do
    b[#b + 1] = tonumber(i)
  end
  local j = 1
  for i = 1, #a do
    while j <= #b and a[i] > b[j] do
      j = j + 1
    end
    if j <= #b and a[i] == b[j] then
      if first then
        io.write(a[i])
        first = false
      else
        io.write("," .. a[i])
      end
      j = j + 1
    end
  end
  print()
end
