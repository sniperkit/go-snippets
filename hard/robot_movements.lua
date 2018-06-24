function robot(f, x, y)
  if x == 4 and y == 4 then return 1 end
  local ret = 0
  if x > 1 and not f[x-1 + 4*(y-1)] then
    local g = {}
    for k, v in pairs(f) do g[k] = f[k] end
    g[x-1 + 4*(y-1)] = true
    ret = ret + robot(g, x-1, y)
  end
  if y > 1 and not f[x + 4*(y-2)] then
    local g = {}
    for k, v in pairs(f) do g[k] = f[k] end
    g[x + 4*(y-2)] = true
    ret = ret + robot(g, x, y-1)
  end
  if x < 4 and not f[x+1 + 4*(y-1)] then
    local g = {}
    for k, v in pairs(f) do g[k] = f[k] end
    g[x+1 + 4*(y-1)] = true
    ret = ret + robot(g, x+1, y)
  end
  if y < 4 and not f[x + 4*y] then
    local g = {}
    for k, v in pairs(f) do g[k] = f[k] end
    g[x + 4*y] = true
    ret = ret + robot(g, x, y+1)
  end
  return ret
end

print(robot({true}, 1, 1))
