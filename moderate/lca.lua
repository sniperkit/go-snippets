parent = {[3]=8, [8]=30, [10]=20, [20]=8, [29]=20, [30]=nil, [52]=30}

for line in io.lines(arg[1]) do
  local n1, n2 = line:match("(%d+) (%d+)")
  n1, n2 = tonumber(n1), tonumber(n2)
  local found = false
  while n1 and n1 ~= n2 do
    local i = n2
    while i do
      if n1 == i then
        found = true
        break
      end
      i = parent[i]
    end
    if found then break end
    n1 = parent[n1]
  end
  print(n1)
end
