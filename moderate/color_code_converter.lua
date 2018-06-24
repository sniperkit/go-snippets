function rint(a)
  return math.floor(a*255 + 0.5)
end

for line in io.lines(arg[1]) do
  local r, g, b = 0, 0, 0
  if line:sub(1, 1) == "#" then
    r, g, b = line:match("(%x%x)(%x%x)(%x%x)")
    r, g, b = tonumber("0x" .. r), tonumber("0x" .. g), tonumber("0x" .. b)
  elseif line:sub(1, 3) == "HSL" then
    h, s, l = line:match("HSL%((%d+),(%d+),(%d+)%)")
    s, l = s/100, l/100
    local c, h2 = (1-math.abs(2*l-1))*s, h/60
    local x, m = c*(1-math.abs(h2%2 - 1)), l-c/2
    if h2 < 1 then
      r, g, b = rint(c+m), rint(x+m), rint(m)
    elseif h2 < 2 then
      r, g, b = rint(x+m), rint(c+m), rint(m)
    elseif h2 < 3 then
      r, g, b = rint(m), rint(c+m), rint(x+m)
    elseif h2 < 4 then
      r, g, b = rint(m), rint(x+m), rint(c+m)
    elseif h2 < 5 then
      r, g, b = rint(x+m), rint(m), rint(c+m)
    elseif h2 < 6 then
      r, g, b = rint(c+m), rint(m), rint(x+m)
    end
  elseif line:sub(1, 3) == "HSV" then
    h, s, v = line:match("HSV%((%d+),(%d+),(%d+)%)")
    s, v = s/100, v/100
    local c, h2 = v*s, h/60
    local x, m = c*(1-math.abs(h2%2 - 1)), v-c
    if h2 < 1 then
      r, g, b = rint(c+m), rint(x+m), rint(m)
    elseif h2 < 2 then
      r, g, b = rint(x+m), rint(c+m), rint(m)
    elseif h2 < 3 then
      r, g, b = rint(m), rint(c+m), rint(x+m)
    elseif h2 < 4 then
      r, g, b = rint(m), rint(x+m), rint(c+m)
    elseif h2 < 5 then
      r, g, b = rint(x+m), rint(m), rint(c+m)
    elseif h2 < 6 then
      r, g, b = rint(c+m), rint(m), rint(x+m)
    end
  elseif line:sub(1, 1) == "(" then
    c, m, y, k = line:match("%(([%d.]+),([%d.]+),([%d.]+),([%d.]+)%)")
    r, g, b = rint((1-c)*(1-k)), rint((1-m)*(1-k)), rint((1-y)*(1-k))
  end
  print("RGB(" .. r .. "," .. g .. "," .. b .. ")")
end
