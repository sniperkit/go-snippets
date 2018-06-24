function pol(o, a, b)
  return (o == '*') and a * b or (o == '/') and a / b or a + b
end

for line in io.lines(arg[1]) do
  local op, num = {}, {}
  for i in line:gmatch("[*/+]") do op[#op + 1] = i end
  for i in line:gmatch("%d+") do num[#num + 1] = i end
  local res = num[1]
  for i = #op, 1, -1 do
    res = pol(op[i], res, num[#num - i + 1])
  end
  print(math.floor(res + 0.001))
end
