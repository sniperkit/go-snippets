for line in io.lines(arg[1]) do
  local sep = line:find(",")
  if sep then
    print(line:sub(1, sep-1):find(line:sub(sep+1) .. "$") and 1 or 0)
  end
end
