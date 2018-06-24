for line in io.lines(arg[1]) do
  local f, b, n = line:match("(%d+) (%d+) (%d+)")
  for i = 1, n do
    if i > 1 then
      io.write(" ")
    end
    if math.fmod(i, f) > 0 and math.fmod(i, b) > 0 then
      io.write(i)
    else
      if math.fmod(i, f) == 0 then
        io.write("F")
      end
      if math.fmod(i, b) == 0 then
        io.write("B")
      end
    end
  end
  print()
end
