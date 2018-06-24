for line in io.lines(arg[1]) do
  local x1, y1, x2, y2 = line:match("(-?%d+) (-?%d+) (-?%d+) (-?%d+)")
  x1, y1, x2, y2 = tonumber(x1), tonumber(y1), tonumber(x2), tonumber(y2)
  print((x1 == x2) and ((y1 == y2) and "here" or
                        (y1 < y2) and "N" or "S") or
        (y1 == y2) and ((x1 < x2) and "E" or "W") or
        ((x1 < x2) and ((y1 < y2) and "NE" or "SE") or
                       ((y1 < y2) and "NW" or "SW")))
end
