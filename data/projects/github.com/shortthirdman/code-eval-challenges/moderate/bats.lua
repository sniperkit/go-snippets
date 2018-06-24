for line in io.lines(arg[1]) do
  local a = {}
  for i in line:gmatch("%d+") do a[#a+1] = tonumber(i) end
  local l, d, n, count = a[1], a[2], a[3], 0
  local t, tx, i = 0, 6 - d, 6
  while i <= l-6 do
    if i > tx-d then
      i = tx
      if t == n then
        tx = l-6+d
      else
        tx = a[t+4]
        t = t + 1
      end
    else
      count = count + 1
    end
    i = i + d
  end
  print(count)
end
