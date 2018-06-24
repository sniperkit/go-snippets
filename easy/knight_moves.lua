function pos(l, lx, r, rx)
  return string.char(string.byte(l) + lx) .. string.char(string.byte(r) + rx)
end

for line in io.lines(arg[1]) do
  local l, r, m = line:sub(1, 1), line:sub(2, 2), {}
  if (l > "b" and r > "1") then
    m[#m+1] = pos(l, -2, r, -1)
  end
  if (l > "b" and r < "8") then
    m[#m+1] = pos(l, -2, r, 1)
  end
  if (l > "a" and r > "2") then
    m[#m+1] = pos(l, -1, r, -2)
  end
  if (l > "a" and r < "7") then
    m[#m+1] = pos(l, -1, r, 2)
  end
  if (l < "h" and r > "2") then
    m[#m+1] = pos(l, 1, r, -2)
  end
  if (l < "h" and r < "7") then
    m[#m+1] = pos(l, 1, r, 2)
  end
  if (l < "g" and r > "1") then
    m[#m+1] = pos(l, 2, r, -1)
  end
  if (l < "g" and r < "8") then
    m[#m+1] = pos(l, 2, r, 1)
  end
  print(table.concat(m, ' '))
end
