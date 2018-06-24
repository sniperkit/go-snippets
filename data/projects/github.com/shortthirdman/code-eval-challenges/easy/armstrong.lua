function div(a, b) return math.floor(a/b) end

for line in io.lines(arg[1]) do
  local a = tonumber(line)
  local e, r = math.floor(math.log10(a)+1), a
  while a > 0 do
    r, a = r - math.pow(a % 10, e), div(a, 10)
  end
  print((r == 0) and "True" or "False")
end
