for line in io.lines(arg[1]) do
  local r = 0
  for i in line:gmatch(".") do
    r = r + i
  end
  print(r)
end
