local shrink = 1.25
for line in io.lines(arg[1]) do
  local a, n = {}, 0
  for i in line:gmatch("%d+") do
    a[#a+1] = tonumber(i)
  end
  local gap, last, changed = #a, 0, false
  while gap > 1 or changed do
    n, changed = n + 1, false
    if gap > 1 then
      gap = math.floor(gap / shrink)
    end
    for j = 1, #a - gap do
      if a[j] > a[j + gap] then
        a[j], a[j + gap], last, changed = a[j + gap], a[j], n, true
      end
    end
  end
  print(last)
end
