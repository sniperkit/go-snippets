m, n = {}, 256
for line in io.lines(arg[1]) do
  if line:find("Query") then
    local q, sum = line:match(" (%d+)"), 0
    if line:find("Row") then
      for i = 1, n do
        if m[(q-1)*n + i] then
          sum = sum + m[(q-1)*n + i]
        end
      end
    else
      for i = 1, n do
        if m[(i-1)*n + q] then
          sum = sum + m[(i-1)*n + q]
        end
      end
    end
    print(sum)
  else
    local p, v = line:match(" (%d+) (%d+)")
    if line:find("Row") then
      for i = 1, n do m[(p-1)*n + i] = v end
    else
      for i = 1, n do m[(i-1)*n + p] = v end
    end
  end
end
