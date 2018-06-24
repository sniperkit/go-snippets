for line in io.lines(arg[1]) do
  if #line > 0 then
    local sep, m = line:find(";"), {}
    for i in line:sub(sep):gmatch("%d+") do
      if m[i] then
        print(i)
        break
      end
      m[i] = true
    end
  end
end
