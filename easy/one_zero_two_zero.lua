function div(a, b) return math.floor(a/b) end

function xz(n, a)
  while a > 0 do
    if a % 2 == 0 then
      if n > 0 then
        n = n - 1
      else
        return false
      end
    end
    a = div(a, 2)
  end
  return n == 0
end

for line in io.lines(arg[1]) do
  local n, a = line:match("(%d+) (%d+)")
  n, a = tonumber(n), tonumber(a)
  local r = 0
  for i = 2 ^ n, a do
    if xz(n, i) then
      r = r + 1
    end
  end
  print(r)
end
