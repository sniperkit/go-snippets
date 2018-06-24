function div(a, b) return math.floor(a/b) end

for line in io.lines(arg[1]) do
  s, n, k = {}, 0, 0
  for i in line:gmatch("[^|]") do
    n = n + 1
    s[#s+1] = string.byte(i) * 2048 + n
    if i == "$" then
      k = n
    end
  end
  table.sort(s)
  for _ = 1, n do
    io.write(string.char(div(s[k], 2048)))
    k = s[k] % 2048
  end
  print()
end
