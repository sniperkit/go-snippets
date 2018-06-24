for line in io.lines(arg[1]) do
  if #line > 0 then
    s = ""
    for i in string.gmatch(line, "%S+") do
      s = i .. " " .. s
    end
    print(s:sub(1, -2))
  end
end
