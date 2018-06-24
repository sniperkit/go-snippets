for line in io.lines(arg[1]) do
  local p, count = {}, 0
  for i in line:gmatch("-?%d+") do
    p[#p+1] = tonumber(i)
  end
  for i = 1, #p/2 - 2 do
    for j = i+1, #p/2 - 1 do
      local dx, dy = p[i*2-1] - p[j*2-1], p[i*2] - p[j*2]
      for k = 1, #p/2 do
        if k ~= i and k ~= j and dx*(p[i*2]-p[k*2]) == (p[i*2-1]-p[k*2-1])*dy then
          if k > j then
            count = count + 1
          end
          break
        end
      end
    end
  end
  print(count)
end
