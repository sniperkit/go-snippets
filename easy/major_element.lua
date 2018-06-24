for line in io.lines(arg[1]) do
  local el, found, n = {}, false, 0
  for i in line:gmatch("%d+") do
    i = tonumber(i)
    el[i], n = el[i] and el[i] + 1 or 1, n + 1
  end
  for k, v in pairs(el) do
    if v > n/2 then
      found = true
      print(k)
      break
    end
  end
  if not found then print("None") end
end
