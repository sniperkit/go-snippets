for line in io.lines(arg[1]) do
  local r = {}
  for i in line:gmatch("%a+") do
    r[#r+1] = i:gsub("%u", string.lower)
  end
  print(table.concat(r, " "))
end
