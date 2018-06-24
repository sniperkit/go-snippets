for line in io.lines(arg[1]) do
  local a, r = {}, 0
  for i = 1, 26 do a[#a+1] = 0 end
  for i in line:gmatch("%a") do
    local c = string.byte(string.lower(i))-96
    a[c] = a[c] +1
  end
  table.sort(a)
  for i = 26, 1, -1 do
    if a[i] == 0 then break end
    r = r + i * a[i]
  end
  print(r)
end
