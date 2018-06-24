for line in io.lines(arg[1]) do
  if #line > 0 then
    local c = line:sub(-1)
    local l = line:find(c .. "[^" .. c .. "]*,")
    print(l and l - 1 or -1)
  end
end
