function div(a, b) return math.floor(a/b) end

for line in io.lines(arg[1]) do
  if #line > 0 then
    local a = tonumber(line)
    if a == 0 then
      print(0)
    else
      local b, c = 0, a
      while a > 0 do
        b, a = b * 2 + a % 2, div(a, 2)
      end
      while c > 0 do
        io.write(b % 2)
        b, c = div(b, 2), div(c, 2)
      end
      print()
    end
  end
end
