for line in io.lines(arg[1]) do
  local count, i = 0, 1
  while true do
    i = line:find(">>%-%->", i)
    if i == nil then break end
    count, i = count + 1, i + 4
  end
  i = 1
  while true do
    i = line:find("<%-%-<<", i)
    if i == nil then break end
    count, i = count + 1, i + 4
  end
  print(count)
end
