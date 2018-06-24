for line in io.lines(arg[1]) do
  local n, c, r = {}, 1, 0
  for i in line:gmatch("%S+") do
    local ix = tonumber(i)
    n[ix], c = n[ix] and 0 or c, c + 1
  end
  for i = 1, 9 do
    if n[i] and n[i] > 0 then
      r = n[i]
      break
    end
  end
  print(r)
end
