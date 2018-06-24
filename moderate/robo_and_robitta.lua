for line in io.lines(arg[1]) do
  local r, a, b, x, y = 0, line:match("(%d+)x(%d+) | (%d+) (%d+)")
  while b ~= y do
    r, a, b, x, y = r+a, b-1, a, b-y, x
  end
  print(r+x)
end
