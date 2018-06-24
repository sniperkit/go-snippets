function justi(s)
  if #s == 80 then
    print(s)
  else
    local r = {}
    for i in s:gmatch("%S+") do
      r[#r+1] = i
    end
    local l = 80 - #table.concat(r)
    for i = 1, l%(#r-1) do
      r[i] = r[i] .. " "
    end
    print(table.concat(r, string.rep(" ", math.floor(l/(#r-1)))))
  end
end

function maxwords(s)
  local t = (#s > 81) and s:sub(1,81) or s
  return t:find(" %S*$") - 1
end

function justify(s)
  local t = s:match"^%s*(.*)"
  if #t <= 80 then
    print(t)
  else
    local m = maxwords(t)
    justi(t:sub(1, m))
    return justify(t:sub(m+1))
  end
end

for line in io.lines(arg[1]) do
  if #line > 80 then
    justify(line)
  else
    print(line)
  end
end
