for line in io.lines(arg[1]) do
  if #line > 0 then
    local sum = 0
    for i in line:gmatch('"id": %d+,') do
      sum = sum + i:match("%d+")
    end
    print(sum)
  end
end
