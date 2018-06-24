for line in io.lines(arg[1]) do
  local m = 10
  for i in line:gmatch("[^,]+") do
    local n = #i:match("X%.*Y") - 2
    if n < m then
      m = n
    end
  end
  print(m)
end
