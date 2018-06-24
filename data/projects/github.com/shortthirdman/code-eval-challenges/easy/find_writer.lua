for line in io.lines(arg[1]) do
  local sep = line:find("|")
  if sep then
    for i in line:sub(sep+1):gmatch("%S+") do
      local ix = tonumber(i)
      io.write(line:sub(ix, ix))
    end
    print()
  end
end
