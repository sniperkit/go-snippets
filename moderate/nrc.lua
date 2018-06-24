for line in io.lines(arg[1]) do
  local a = {}
  for i in line:gmatch(".") do
    a[i] = a[i] and a[i] + 1 or 1
  end
  for i in line:gmatch(".") do
    if a[i] == 1 then
      io.write(i)
      break
    end
  end
  print()
end
