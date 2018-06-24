function happy (a)
  local ret = 0
  while #a > 0 do
    local b = tonumber(a:sub(#a))
    ret, a = ret +  b * b, a:sub(1, #a-1)
  end
  return tostring(ret)
end

for line in io.lines(arg[1]) do
  local a, b = line, {line}
  for _ = 1, 127 do
    if a == "1" or a == "0" then
      break
    end
    a = happy(a)
    for _, j in ipairs(b) do
      if j == a then
        a = "0"
        break
      end
    end
    b[#b + 1] = a
  end
  print((a == "1") and 1 or 0)
end
