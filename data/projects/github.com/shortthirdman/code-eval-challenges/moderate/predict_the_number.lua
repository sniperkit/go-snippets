function div(a, b) return math.floor(a/b) end

for line in io.lines(arg[1]) do
  local n, i = tonumber(line), 0
  while n > 0 do
    if n%2 > 0 then i = i + 1 end
    n = div(n, 2)
  end
  print(i%3)
end
