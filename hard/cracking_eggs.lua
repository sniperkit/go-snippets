ht, ec = {}, 2^16
function floors(e, s)
  if e == 0 or s == 0 then
    return 0
  elseif e == 1 then
    return s
  elseif s == 1 then
    return 1
  elseif e > s then
    return floors(s, s)
  end
  if not ht[e*ec+s] then
    ht[e*ec+s] = floors(e-1, s-1) + floors(e, s-1) + 1
  end
  return ht[e*ec+s]
end

for line in io.lines(arg[1]) do
  local e, s = line:match("(%d+) (%d+)")
  e, s, n = tonumber(e), tonumber(s), 0
  while floors(e, n) < s do
    n = n + 1
  end
  print(n)
end
