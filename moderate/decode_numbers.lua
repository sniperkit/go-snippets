function decode(s)
  while true do
    if #s < 2 then return 1 end
    if tonumber(s:sub(1,1)) > 2 or s:sub(1,1) == "0" then
      s = s:sub(2)
    elseif s:sub(2,2) == "0" or (s:sub(1,1) == "2" and tonumber(s:sub(2,2)) > 6) then
      s = s:sub(3)
    else
      if #s == 2 then return 2 end
      return decode(s:sub(2)) + decode(s:sub(3))
    end
  end
end

for line in io.lines(arg[1]) do
  print(decode(line))
end
