function div(a, b) return math.floor(a/b) end

for line in io.lines(arg[1]) do
  local sep, a = line:find(";"), {}
  for i in line:sub(1, sep):gmatch("%d+") do
    a[#a + 1] = i
  end
  local l = tonumber(line:sub(sep+1))
  for c = l, #a, l do
    for i = 1, div(l, 2) do
      a[c+1-i], a[c-l+i] = a[c-l+i], a[c+1-i]
    end
  end
  print(table.concat(a, ","))
end
