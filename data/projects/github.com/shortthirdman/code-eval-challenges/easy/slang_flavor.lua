local phrases = {", yeah!", ", this is crazy, I tell ya.", ", can U believe this?", ", eh?", ", aw yea.", ", yo.", "? No way!", ". Awesome!"}

local i, l = 1, false
for line in io.lines(arg[1]) do
  for c in line:gmatch(".") do
    if c == '.' or c == '!' or c == '?' then
      if l then
        io.write(phrases[i])
        l, i = false, i % 8 + 1
      else
        io.write(c)
        l = true
      end
    else
      io.write(c)
    end
  end
  print()
end
