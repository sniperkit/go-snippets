local n, i, b = 0, 0, {}
for line in io.lines(arg[1]) do
  if n == 0 then
    n = tonumber(line)
    i = n
  else
    if i == n then
      local c = 1
      for j in string.gmatch(line, "[^,]+") do
        if c == 1 then
          b[1] = j
        else
          b[c] = b[c-1] + j
        end
        c = c + 1
      end
    else
      local c = 1
      for j in string.gmatch(line, "[^,]+") do
        if c == 1 then
          b[1] = b[1] + j
        else
          b[c] = math.min(b[c-1], b[c]) + j
        end
        c = c + 1
      end
    end
    i = i - 1
    if i == 0 then
      print(b[n])
      n = 0
    end
  end
end
