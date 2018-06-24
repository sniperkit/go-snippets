for line in io.lines(arg[1]) do
  h, b, r = {}, {}, {}
  s1, s2, s3, s4 = line:match("%[([-]?%d+),([-]?%d+)%] %[([-]?%d+),([-]?%d+)%]")
  h[1], h[2] = math.abs(tonumber(s1)-tonumber(s3)), math.abs(tonumber(s2)-tonumber(s4))
  table.sort(h)
  for i in line:gmatch("%d+ %[[-]?%d+,[-]?%d+,[-]?%d+%] %[[-]?%d+,[-]?%d+,[-]?%d+%]") do
    ix, s1, s2, s3, s4, s5, s6 = i:match("(%d+) %[([-]?%d+),([-]?%d+),([-]?%d+)%] %[([-]?%d+),([-]?%d+),([-]?%d+)%]")
    b[1], b[2], b[3] = math.abs(tonumber(s1)-tonumber(s4)), math.abs(tonumber(s2)-tonumber(s5)), math.abs(tonumber(s3)-tonumber(s6))
    table.sort(b)
    if b[1] <= h[1] and b[2] <= h[2] then -- or diagonal(b, h)
      r[#r + 1] = tonumber(ix)
    end
  end
  if #r > 0 then
    table.sort(r)
    print(table.concat(r, ","))
  else
    print("-")
  end
end
