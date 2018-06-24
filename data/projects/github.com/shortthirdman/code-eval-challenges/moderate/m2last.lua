for line in io.lines(arg[1]) do
  local b = {}
  for i in string.gmatch(line, "%S+") do
    b[#b + 1] = i
  end
  local i = #b - b[#b]
  if i > 0 then
    print(b[i])
  end
end
