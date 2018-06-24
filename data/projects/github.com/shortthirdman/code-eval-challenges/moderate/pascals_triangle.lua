for line in io.lines(arg[1]) do
  io.write(1)
  for i = 1, line-1 do
    io.write(" 1")
    local r = 1
    for j = 1, i do
      r = (r*(i+1-j))/j
      io.write(" " .. r)
    end
  end
  print()
end
