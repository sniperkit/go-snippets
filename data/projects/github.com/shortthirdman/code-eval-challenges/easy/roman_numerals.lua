ronum = {1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
rostr = {"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
for line in io.lines(arg[1]) do
  local i, num = 1, tonumber(line)
  while num > 0 do
    if num >= ronum[i] then
      io.write(rostr[i])
      num = num - ronum[i]
    else
      i = i + 1
    end
  end
  print()
end
