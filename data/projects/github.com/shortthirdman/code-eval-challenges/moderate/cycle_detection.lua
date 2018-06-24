function uniq(s)
  for i = 2, #s do
    if s[i] == s[1] then return false end
  end
  return true
end

for line in io.lines(arg[1]) do
  local s = {}
  for i in line:gmatch("%S+") do s[#s + 1] = i end
  while #s > 1 and uniq(s) do
    table.remove(s, 1)
  end
  if #s == 1 then
    print()
  else
    local t = {s[1]}
    for i = 2, #s do
      if s[i] == s[1] then break end
      t[#t + 1] = s[i]
    end
    print(table.concat(t, " "))
  end
end
