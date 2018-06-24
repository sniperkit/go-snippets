function div(a, b) return math.floor(a/b) end

t = {}
function pali(a)
  local y, c = 1, a
  while c > 0 do
    t[y], y, c = c % 10, y + 1, div(c, 10)
  end
  y = y - 1
  for x = 1, y/2 do
    if t[x] ~= t[y] then
      return false
    end
    y = y - 1
  end
  return true
end

for line in io.lines(arg[1]) do
  local l, r = line:match("(%d+) (%d+)")
  local n = 0
  for i = l, r do
    local prev = -1
    for j = i, r do
      local p = 0
      if prev > -1 then
        p = prev
        if pali(j) then p = p + 1 end
      else
        p = 0
        for k = i, j do
          if pali(k) then p = p + 1 end
        end
      end
      if p%2 == 0 then n = n + 1 end
      prev = p
    end
  end
  print(n)
end
