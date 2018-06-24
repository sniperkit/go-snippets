for line in io.lines(arg[1]) do
  local l, m = 0, -math.huge
  for i in line:gmatch("[-]?%d+") do
    a = tonumber(i)
    if a+l > m then m = a+l end
    if a+l > 0 then l = l+a else l = 0 end
  end
  print(m)
end
