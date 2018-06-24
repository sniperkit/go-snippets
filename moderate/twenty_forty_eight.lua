function slide(a)
  local ln, li, lz = nil, nil, nil
  for i = 1, n do
    if a[i] == 0 then
      lz = lz or i
    else
      if ln == a[i] then
        a[li], ln, a[i] = 2 * ln, nil, 0
        lz = lz or i
      else
        ln = a[i]
        if not lz then
          li = i
        else
          li, a[lz], a[i], lz = lz, a[i], 0, lz + 1
        end
      end
    end
  end
  return a
end

for line in io.lines(arg[1]) do
  local sep = line:find("%d;") + 3
  n = line:match("(%d+);")
  local b, i = {}, 0
  for s in line:sub(sep):gmatch("%d+") do
    if i % n == 0 then
      b[#b+1] = {}
    end
    i = i + 1
    b[#b][#b[#b]+1] = tonumber(s)
  end

  for i = 1, n do
    local t = {}
    if line:sub(1,1) == 'L' then
      t = b[i]
      t = slide(t)
      b[i] = t
    elseif line:sub(1,1) == 'R' then
      for j = 1, n do
        t[j] = b[i][n-j+1]
      end
      t = slide(t)
      for j = 1, n do
        b[i][j] = t[n-j+1]
      end
    elseif line:sub(1,1) == 'U' then
      for j = 1, n do
        t[j] = b[j][i]
      end
      t = slide(t)
      for j = 1, n do
        b[j][i] = t[j]
      end
    elseif line:sub(1,1) == 'D' then
      for j = 1, n do
        t[j] = b[n-j+1][i]
      end
      t = slide(t)
      for j = 1, n do
        b[j][i] = t[n-j+1]
      end
    end
  end

  local r = {}
  for i = 1, n do
    r[i] = table.concat(b[i], " ")
  end
  print(table.concat(r, "|"))
end
