for line in io.lines(arg[1]) do
  local a, swap = {}, false
  for i in line:gmatch("%S+") do
    if swap then
      local x, y = i:match("(%d+)-(%d+)")
      a[x+1], a[y+1] = a[y+1], a[x+1]
    elseif i == ":" then
      swap = true
    else
      a[#a + 1] = i
    end
  end
  print(table.concat(a, " "))
end
