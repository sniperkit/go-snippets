for line in io.lines(arg[1]) do
  local up = true
  for i in line:gmatch(".") do
    if i:find("%a") then
      if up then
        io.write(i:upper())
        up = false
      else
        io.write(i:lower())
        up = true
      end
    else
      io.write(i)
    end
  end
  print()
end
