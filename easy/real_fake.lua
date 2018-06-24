for line in io.lines(arg[1]) do
  local a, su = {}, 0
  for i in line:gmatch("%d") do
    a[#a+1] = tonumber(i)
  end
  for i = #a-1, 1, -2 do
    a[i] = 2*a[i]
  end
  for i = 1, #a do
    su = su + a[i]
  end
  print((su%10 == 0) and "Real" or "Fake")
end
