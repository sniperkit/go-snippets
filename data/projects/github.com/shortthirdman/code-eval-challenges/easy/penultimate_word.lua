for line in io.lines(arg[1]) do
  for i in line:gmatch("%S+") do
    c, p = i, c
  end
  print(p and p or "")
end
