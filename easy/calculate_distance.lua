for line in io.lines(arg[1]) do
  local x1, y1, x2, y2 = line:match("%((-?%d+), (-?%d+)%) %((-?%d+), (-?%d+)%)")
  local x, y = tonumber(x1) - tonumber(x2), tonumber(y1) - tonumber(y2)
  print(math.floor(math.sqrt(x * x + y * y)))
end
