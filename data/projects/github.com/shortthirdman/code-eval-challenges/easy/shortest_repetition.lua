for line in io.lines(arg[1]) do
  local p = 1
  for i = 2, #line do
    if line:sub(i, i) ~= line:sub(i-p, i-p) then
      p = i
    end
  end
  print(p)
end
