function subs(s, t)
  if #t == 0 then return 1 end
  while #s > 0 do
    if s:sub(1, 1) == t:sub(1, 1) then
      return subs(s:sub(2), t:sub(2)) + subs(s:sub(2), t)
    end
    s = s:sub(2)
  end
  return 0
end

for line in io.lines(arg[1]) do
  print(subs(line:match("([^,]+),([^,]+)")))
end
