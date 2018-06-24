function div(a, b) return math.floor(a/b) end

sq = {[4]=2, [9]=3}
function set(a, b)
  if a % (b + b) < b then
    return a + b, true
  end
  return a, false
end

for line in io.lines(arg[1]) do
  local t, valid, n = {}, true, tonumber(line:match("(%d);"))
  for i in line:sub(3):gmatch("%d") do
    t[#t+1] = tonumber(i)
  end
  local csq, col = {}, {}
  for i = 0, sq[n]-1 do csq[i] = 0 end
  for i = 1, n do col[i] = 0 end
  for i = 1, n do
    local crow = 0
    for j = 1, n do
      d = 2^(t[(i-1)*n+j]-1)
      crow, ok = set(crow, d)
      if not ok then valid = false; break end
      csq[div(j-1, sq[n])], ok = set(csq[div(j-1, sq[n])], d)
      if not ok then valid = false; break end
      col[j], ok = set(col[j], d)
      if not ok then valid = false; break end
    end
    if i % sq[n] == 0 then
      for j = 0, sq[n]-1 do csq[j] = 0 end
    end
    if not valid then break end
  end
  print(valid and "True" or "False")
end
