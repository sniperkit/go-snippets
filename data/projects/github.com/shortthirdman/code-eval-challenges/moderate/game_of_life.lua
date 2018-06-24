local n, l, c, d = 0, 0, {}, {}
for line in io.lines(arg[1]) do
  if n == 0 then
    n = #line
    d = {-n-1, -n, 1-n, -1, 0, 1, n-1, n, n+1}
  end
  local i = 0
  while true do
    i = line:find("*", i+1)
    if not i then break end
    c[#c+1] = l*n + i - 1
  end
  l = l + 1
end

for _ = 1, 10 do
  m = {}
  for _, i in ipairs(c) do
    for jx, j in ipairs(d) do
      if j ~= 0 and i+j >= 0 and i+j < n*n and (i%n > 0 or jx%3 ~= 1) and (i%n < n-1 or jx%3 > 0) then
        m[i+j] = m[i+j] and m[i+j] + 1 or 1
      end
    end
  end
  local c2 = {}
  l = 1
  for i = 0, n*n-1 do
    if m[i] == 2 then
      while l <= #c and c[l] < i do l = l + 1 end
      if l <= #c and c[l] == i then
        c2[#c2 + 1] = i
      end
    elseif m[i] == 3 then
      c2[#c2 + 1] = i
    end
  end
  c = c2
end

l = 1
for i = 0, n-1 do
  for j = 0, n-1 do
    if #c >= l and i*n+j == c[l] then
      io.write("*")
      l = l + 1
    else
      io.write(".")
    end
  end
  print()
end
