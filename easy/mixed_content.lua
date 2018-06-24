for line in io.lines(arg[1]) do
  local s, n = {}, {}
  for i in line:gmatch("[^,]+") do
    if i:find("^%d+$") then
      n[#n+1] = i
    else
      s[#s+1] = i
    end
  end
  if #s > 0 then
    io.write(s[1])
    for i = 2, #s do
      io.write("," .. s[i])
    end
    if #n > 0 then io.write("|") end
  end
  if #n > 0 then
    io.write(n[1])
    for i = 2, #n do
      io.write("," .. n[i])
    end
  end
  print()
end
