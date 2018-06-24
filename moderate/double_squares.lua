sq, n = {}, 0
for i = 0, 46340 do sq[i*i] = true end
for line in io.lines(arg[1]) do
  if n == 0 then
    n = line
  else
    local x = tonumber(line)
    local num, l, bot, top = 0, {}, math.floor(math.sqrt(x/2)), math.floor(math.sqrt(x))
    for i = bot, top do
      local t = x - i*i
      if sq[t] and not l[t] then
        l[i*i], num = true, num + 1
      end
    end
    print(num)
  end
end
