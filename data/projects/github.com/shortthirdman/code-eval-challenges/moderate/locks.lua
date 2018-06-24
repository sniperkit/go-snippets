for line in io.lines(arg[1]) do
  local n, m = line:match("(%d+) (%d+)")
  if m == "0" then
    print(n)
  elseif m == "1" then
    print(n-1)
  else
    n, m = tonumber(n), tonumber(m)
    local l, count = {}, 0
    for _ = 0, m%2 do
      for j = 1, n-1, 2 do l[j] = true end
      for j = 2, n-1, 3 do
        if l[j] then l[j] = false else l[j] = true end
      end
    end
    if l[n-1] then l[n-1] = false else l[n-1] = true end
    for i = 0, n-1 do
      if not l[i] then count = count + 1 end
    end
    print(count)
  end
end
