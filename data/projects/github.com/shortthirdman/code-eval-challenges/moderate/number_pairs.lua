for line in io.lines(arg[1]) do
  local sep, v, vk, res = line:find(";"), {}, {}, {}
  local n = tonumber(line:sub(sep+1))
  for i in line:sub(1, sep):gmatch("%d+") do
    v[tonumber(i)], vk[#vk+1] = true, i
  end
  for i = 1, #vk do
    if 2*vk[i] < n and v[n-vk[i]] then
      res[#res+1] = vk[i] .. "," .. tostring(n-vk[i])
    end
  end
  print((#res > 0) and table.concat(res, ";") or "NULL")
end
