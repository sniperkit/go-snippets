for line in io.lines(arg[1]) do
  local sep = line:find(" ")
  for i = 1, sep-1 do
    io.write((line:sub(sep+i, sep+i) == "1") and line:sub(i, i):upper() or line:sub(i, i))
  end
  print()
end
