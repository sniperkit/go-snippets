for line in io.lines(arg[1]) do
  line = line:gsub("[^0-9a-j]", "")
  for i in line:gmatch(".") do
    io.write(i:match("[a-j]") and i:byte()-string.byte("a") or i)
  end
  print((#line > 0) and "" or "NONE")
end
