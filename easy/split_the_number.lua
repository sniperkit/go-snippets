for line in io.lines(arg[1]) do
  local sep = line:find(" ")
  if sep then
    local c, o, r, v = 1, 1, 0, 0
    for i in line:sub(sep+1):gmatch(".") do
      if i:find("%a") then
        v = v*10 + tonumber(line:sub(c, c))
        c = c + 1
      elseif i == "+" then
        o, r, v = 1, r + v*o, 0
      elseif i == "-" then
        o, r, v = -1, r + v*o, 0
      end
    end
    print(r + v*o)
  end
end
