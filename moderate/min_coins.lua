function div(a, b) return math.floor(a/b) end

b = {[0]=0, 1, 2, 1, 2}
for line in io.lines(arg[1]) do
  print(div(line, 5) + b[line%5])
end
