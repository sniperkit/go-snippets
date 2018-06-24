for line in io.lines(arg[1]) do
  local a, b, r, p, f = {}, {}, {}, "", false
  for i in line:gmatch("%S+") do
    if f then
      p = {}
      for j in i:gmatch(".") do
        p[#p + 1] = j
      end
      table.sort(p)
    elseif i == "|" then
      f = true
    else
      a[#a + 1] = i
      b[#b + 1] = {}
      for j in i:gmatch(".") do
        b[#b][#b[#b] + 1] = j
      end
      table.sort(b[#b])
    end
  end
  f = false
  for i = 1, #a do
    local q = {}
    for j = 1, #p do
      q[j] = p[j]
    end
    while #q > 0 and #b[i] > 0 do
      if q[#q] < b[i][#b[i]] then
        b[i][#b[i]] = nil
      elseif q[#q] == b[i][#b[i]] then
        b[i][#b[i]] = nil
        q[#q] = nil
      else
        break
      end
    end
    if #q == 0 then
      if f then
        io.write(" ")
      end
      io.write(a[i])
      f = true
    end
  end
  if not f then
    print("False")
  else
    print()
  end
end
