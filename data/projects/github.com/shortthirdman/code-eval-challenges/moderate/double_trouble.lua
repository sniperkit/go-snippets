for line in io.lines(arg[1]) do
  local n, r = #line / 2, 1
  for i = 1, n do
    local a, b = line:sub(i, i), line:sub(i+n, i+n)
    if (a == "A" and b == "B") or
       (a == "B" and b == "A") then
      r = 0
      break
    elseif a == "*" and b == "*" then
      r = r * 2
    end
  end
  print(r)
end
