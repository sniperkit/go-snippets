for line in io.lines(arg[1]) do
  for i in line:gmatch(".") do
    io.write(i:find("[a-z]") and (i:gsub("%l", string.upper)) or
             i:find("[A-Z]") and (i:gsub("%u", string.lower)) or i)
  end
  print()
end
