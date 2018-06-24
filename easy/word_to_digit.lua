a = {zero = 0, one = 1, two = 2, three = 3, four = 4,
     five = 5, six = 6, seven = 7, eight = 8, nine = 9}
for line in io.lines(arg[1]) do
  for i in line:gmatch("[^;]+") do
    io.write(a[i])
  end
  print()
end
