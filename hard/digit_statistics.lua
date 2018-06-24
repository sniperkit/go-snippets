function div(a, b) return math.floor(a/b) end

stats = {{2, 4, 8, 6}, {3, 9, 7, 1}, {4, 6}, {5}, {6},
         {7, 9, 3, 1}, {8, 4, 2, 6}, {9, 1}}
for line in io.lines(arg[1]) do
  a, n = line:match("(%d) (%d+)")
  a = a - 1
  m, res = div(n, #stats[a]), {[0]=0,0,0,0,0,0,0,0,0,0}
  for _, i in ipairs(stats[a]) do res[i] = res[i] + m end
  for i = 1, n % #stats[a] do res[stats[a][i]] = res[stats[a][i]] + 1 end
  for i = 9, 0, -1 do
    res[i+1] = tostring(i) .. ": " .. tostring(res[i])
  end
  print(table.concat(res, ", "))
end
