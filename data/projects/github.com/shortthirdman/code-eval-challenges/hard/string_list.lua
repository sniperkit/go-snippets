function sortu(s)
  local r = {}
  for i in s:gmatch(".") do
    r[#r + 1] = i
  end
  table.sort(r)
  local t = {}
  for i = 1, #r do
    if r[i] ~= r[i-1] then
      t[#t + 1] = r[i]
    end
  end
  return table.concat(t)
end

function slist(s, t, n)
  if n == 0 then return {t} end
  local ret = {}
  for i in s:gmatch(".") do
    for _, j in ipairs(slist(s, t .. i, n-1)) do
      ret[#ret + 1] = j
    end
  end
  return ret
end

for line in io.lines(arg[1]) do
  local n, s = line:match("(%d+),(.+)")
  print(table.concat(slist(sortu(s), "", n), ","))
end
