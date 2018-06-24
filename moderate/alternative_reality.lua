m = {50, 25, 10, 5, 1}

function alter(n, p)
  if n == 0 then return 1 end
  while m[p] > n do p = p + 1 end
  if m[p] == 1 then return 1 end
  return alter(n - m[p], p) + alter(n, p + 1)
end

for line in io.lines(arg[1]) do
  print(alter(tonumber(line), 1))
end
