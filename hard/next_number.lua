function div(a, b) return math.floor(a/b) end

function dsig(a)
  local r = 0
  while a > 0 do
    if a%10 > 0 then
      r = r + 2^(3*(a%10))
    end
    a = div(a, 10)
  end
  return r
end

for line in io.lines(arg[1]) do
  local d = tonumber(line)
  local ds, e = dsig(d), d + 9
  while dsig(e) ~= ds do
    e = e + 9
  end
  print(e)
end
