for line in io.lines(arg[1]) do
  local a = tonumber(line)
  if a < 2 then
    print(a)
  else
    local b, c = 1, 1
    while a > 2 do
      a, b, c = a-1, b+c, b
    end
    print(b)
  end
end
