for line in io.lines(arg[1]) do
  local a = tonumber(line)
  local b = math.floor(a)
  io.write(b .. ".")
  a = (a - b) * 60
  b = math.floor(a)
  io.write(string.format("%02d'", b))
  a = (a - b) * 60
  b = math.floor(a)
  print(string.format("%02d\"", b))
end
