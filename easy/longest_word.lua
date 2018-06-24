for line in io.lines(arg[1]) do
  local m = ""
  for i in string.gmatch(line, "%S+") do
    if #i > #m then m = i end
  end
  print(m)
end
