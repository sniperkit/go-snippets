for line in io.lines(arg[1]) do
  local b, c = "", 1
  for i in string.gmatch(line, "%S+") do
    if i == b then
      c = c + 1
    else
      if #b > 0 then
        io.write(c .. " " .. b .. " ")
        c = 1
      end
      b = i
    end
  end
  print(c .. " " .. b)
end
