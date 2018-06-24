for line in io.lines(arg[1]) do
  local s = {}
  for i in line:gmatch("[-]?%d+") do
    s[#s + 1] = tonumber(i)
  end
  local jolly, n = true, tonumber(s[1])
  if n > 1 then
    local u = {}
    for i = 3, #s do
      local x = math.abs(s[i] - s[i-1])
      if x >= n or x == 0 or u[x-1] == true then
        jolly = false
        break
      end
      u[x-1] = true
    end
  end
  print(jolly and "Jolly" or "Not jolly")
end
