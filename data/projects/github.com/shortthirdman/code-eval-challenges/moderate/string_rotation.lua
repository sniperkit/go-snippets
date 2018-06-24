for line in io.lines(arg[1]) do
  local sep, found = line:find(","), false
  local p, i = 0, line:sub(sep+1):find(line:sub(1, 1))
  while i do
    i = p + i
    if line:sub(sep+i) .. line:sub(sep+1, sep+i-1) == line:sub(1, sep-1) then
      found = true
      break
    end
    p, i = i, line:sub(sep+i+1):find(line:sub(1, 1))
  end
  print(found and "True" or "False")
end
