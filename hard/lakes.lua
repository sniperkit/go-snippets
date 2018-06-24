function div(a, b) return math.floor(a/b) end

for line in io.lines(arg[1]) do
  local r, l, c, n, ix, jx = 0, {}, {}, div(line:find("|"), 2), 1, 1
  for i = 0, n+1 do
    c[i], l[i] = 0, 0
  end
  for i in line:gmatch("[#o]") do
    if i == "o" then
      if l[jx-1] > 0 then
        c[jx] = l[jx-1]
      elseif l[jx] > 0 then
        c[jx] = l[jx]
      elseif c[jx-1] > 0 then
        c[jx] = c[jx-1]
      end
      if l[jx+1] > 0 then
        if c[jx] == 0 then
          c[jx] = l[jx+1]
        elseif c[jx] ~= l[jx+1] then
          local t = l[jx+1]
          for k, v in pairs(l) do
            if v == t then
              l[k] = c[jx]
            end
          end
          for k, v in pairs(c) do
            if v == t then
              c[k] = c[jx]
            end
          end
          r = r - 1
        end
      end
      if c[jx] == 0 then
        c[jx] = ix * 32 + jx
        r = r + 1
      end
    end
    if jx < n then
      jx = jx + 1
    else
      ix = ix + 1
      jx = 1
      for j = 0, n+1 do
        l[j], c[j] = c[j], 0
      end
    end
  end
  print(r)
end
