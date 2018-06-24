for line in io.lines(arg[1]) do
  x1, y1, r, x2, y2 = line:match("Center: %((-?%d+.?%d*), (-?%d+.?%d*)%); Radius: (-?%d+.?%d*); Point: %((-?%d+.?%d*), (-?%d+.?%d*)%)")
  local x, y = x2-x1, y2-y1
  print(x*x + y*y <= r*r)
end
