for line in io.lines(arg[1]) do
  r, m = {}, ""
  for i in string.gmatch(line, "%S+") do
    if #i > #m then
      m = i
    end
  end
  for i = 1, #m do
    r[i] = string.rep("*", i-1) .. m:sub(i, i)
  end
  print(table.concat(r, " "))
end
