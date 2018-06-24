mat = {}
function uniq(x, y, cw, ch)
  local d = {}
  for i = x, x+cw-1 do
    for j = y, y+ch-1 do
      if d[mat[j]:sub(i, i)] then return false end
      d[mat[j]:sub(i, i)] = true
    end
  end
  return true
end

h, w = 0, 0
for line in io.lines(arg[1]) do
  h = h + 1
  if w == 0 then w = #line end
  mat[h] = line
end

res, ma, ch, cw, x, y = {}, 1, 1, 1, 1, 1
while ch <= h do
  if uniq(x, y, cw, ch) then
    if cw*ch > ma then
      ma, res = cw*ch, {{x, y, cw, ch}}
    else
      res[#res + 1] = {x, y, cw, ch}
    end
    while x+cw <= w and uniq(x, y, cw+1, ch) do
      cw = cw + 1
      ma, res = cw*ch, {{x, y, cw, ch}}
    end
  end
  if x+cw <= w then
    x = x + 1
  elseif y+ch <= h then
    x, y = 1, y + 1
  else
    x, y, cw, ch = 1, 1, 1, ch+1
    while cw*ch < ma do cw = cw + 1 end
  end
end

for i = 1, #res do
  for j = res[i][2], res[i][2]+res[i][4]-1 do
    mat[j] = mat[j]:sub(1, res[i][1]-1) .. string.rep("*", res[i][3]) .. mat[j]:sub(res[i][1]+res[i][3])
  end
end
print(table.concat(mat, "\n"))
