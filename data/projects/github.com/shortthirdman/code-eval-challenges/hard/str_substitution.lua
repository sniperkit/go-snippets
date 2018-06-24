for line in io.lines(arg[1]) do
  local s, t, sep, c, f = {}, {}, line:find(";"), 1, true
  for i in line:sub(sep+1):gmatch(".") do
    if f then
      s[c] = s[c] or ""
      if i == ',' then
        f = false
        t[c] = ""
      else
        s[c] = s[c] .. i
      end
    else
      if i == ',' then
        f = true
        c = c + 1
      else
        t[c] = t[c] .. i:gsub("0", "a"):gsub("1", "b")
      end
    end
  end
  local r = line:sub(1, sep-1)
  for i = 1, c do
    if t[i] ~= "" then
      r = r:gsub(s[i], t[i])
    else
      r = r:gsub(s[i], "x")
    end
  end
  r = r:gsub("x", ""):gsub("a", "0"):gsub("b", "1")
  print(r)
end
