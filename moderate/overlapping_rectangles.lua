function within(a, r, p, q)
  return a[p] >= a[r] and a[p] <= a[r+2] and a[q] >= a[r+3] and a[q] <= a[r+1]
end

for line in io.lines(arg[1]) do
  a = {}
  for i in line:gmatch("[-]?%d+") do a[#a + 1] = tonumber(i) end
  print((within(a, 1, 5, 6) or within(a, 1, 7, 8) or within(a, 1, 5, 8) or within(a, 1, 7, 6) or
         within(a, 5, 1, 2) or within(a, 5, 3, 4) or within(a, 5, 1, 4) or within(a, 5, 3, 2)) and
        "True" or "False")
end
