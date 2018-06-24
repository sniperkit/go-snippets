for line in io.lines(arg[1]) do
  if #line > 55 then
    local i, p = 0, 41
    while true do
      i = string.find(line, " ", i+1)
      if not i or i > 40 then break end
      p = i
    end
    line = line:sub(1, p-1) .. "... <Read More>"
  end
  print(line)
end
