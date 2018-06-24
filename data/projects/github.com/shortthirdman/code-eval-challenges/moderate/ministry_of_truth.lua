function truth(a, b)
  local ix = a:find(b)
  if ix then
    return string.rep("_", ix-1) .. b .. string.rep("_", #a-#b-ix+1), true
  end
  return string.rep("_", #a), false
end

for line in io.lines(arg[1]) do
  local sep, t, u, c, res = line:find(";"), {}, {}, 1, {}
  for i in line:sub(sep+1):gmatch("%S+") do u[#u + 1] = i end
  for i in line:sub(1, sep-1):gmatch("%S+") do
    local r, f = "", false
    if c <= #u then
      r, f = truth(i, u[c])
      if f then c = c + 1 end
    else
      r = string.rep("_", #i)
    end
    if #r > 0 then res[#res + 1] = r end
  end
  print((c <= #u) and "I cannot fix history" or table.concat(res, " "))
end
