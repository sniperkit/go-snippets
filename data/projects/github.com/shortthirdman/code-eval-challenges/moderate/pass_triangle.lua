a = {}
for line in io.lines(arg[1]) do
  local s = {}
  for i in line:gmatch("[-]?%d+") do
    s[#s + 1] = tonumber(i)
  end
  if #a > 0 then s[1] = s[1] + a[1] end
  if #a > 1 then s[#s] = s[#s] + a[#a] end
  for i = 2, #s-1 do
    s[i] = s[i] + math.max(a[i-1], a[i])
  end
  a = s
end
table.sort(a)
print(a[#a])
