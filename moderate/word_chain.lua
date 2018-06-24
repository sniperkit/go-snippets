function lwc(f, c)
  local m = #c/2
  if #f == 0 then return m end
  for ix, i in ipairs(f) do
    if #c == 0 or i:sub(-1) == c:sub(1, 1) then
      local g, h = {}, {}
      for jx, j in ipairs(f) do
        h[j:sub(1,1)] = true
      end
      for jx, j in ipairs(f) do
        if ix ~= jx and h[j:sub(-1)] then g[#g + 1] = j end
      end
      local l = lwc(g, i .. c)
      if l > m then
        m = l
      end
    end
  end
  return m
end

for line in io.lines(arg[1]) do
  local s = {}
  for i in string.gmatch(line, "[^,]+") do
    s[#s + 1] = i:sub(1, 1) .. i:sub(-1)
  end
  local l = lwc(s, "")
  print((l == 1) and "None" or l)
end
