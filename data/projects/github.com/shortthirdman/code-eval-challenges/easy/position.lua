for line in io.lines(arg[1]) do
  local x, p1, p2 = line:match("(%d+),(%d+),(%d+)")
  p1, p2 = 2^(p1-1), 2^(p2-1)
  print((x % (p1 + p1) >= p1) == (x % (p2 + p2) >= p2))
end
