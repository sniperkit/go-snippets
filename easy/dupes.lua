for line in io.lines(arg[1]) do
  local p = nil
  for i in line:gmatch("%d+") do
    if p == nil then
      io.write(i)
    elseif p ~= i then
      io.write("," .. i)
    end
    p = i
  end
  print()
end
