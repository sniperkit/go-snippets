function div(a, b) return math.floor(a/b) end

for line in io.lines(arg[1]) do
  n = tonumber(line)
  local m, r = n, 0
  a = {}
  for i = 0, #line-1 do
    a[i] = 0
  end
  self = true
  while m > 0 do
    local v = m % 10
    if v <= #a then
      a[v] = a[v] + 1
    else
      self = false
      break
    end
    m = div(m, 10)
  end
  if self then
    for i = 0, #a do
      r = r*10 + a[i]
    end
    self = n == r
  end
  print(self and 1 or 0)
end
