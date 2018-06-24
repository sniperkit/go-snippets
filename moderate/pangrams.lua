for line in io.lines(arg[1]) do
  f, line = false, line:gsub("%u", string.lower)
  for i = string.byte('a'), string.byte('z') do
    if not line:find(string.char(i)) then
      io.write(string.char(i))
      f = true
    end
  end
  print(f and "" or "NULL")
end
