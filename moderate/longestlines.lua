lol = {}
for line in io.lines(arg[1]) do
  if not n then
    n = line
  else
    for i = 1, n do
      if not lol[i] then
        lol[i] = line
        break
      elseif #line > #lol[i] then
        lol[i], line = line, lol[i]
      end
    end
  end
end
print(table.concat(lol, "\n"))
