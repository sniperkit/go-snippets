for line in io.lines(arg[1]) do
  if not #line then
  elseif not curr then
    curr = line:find("C") or line:find("_")
    print(line:sub(1, curr-1) .. "|" .. line:sub(curr+1))
  else
    px = curr
    local a, b = (px < 2) and 1 or px-1, (px+2 > #line) and #line or px+2
    curr = (line:sub(a, b):find("C") or line:sub(a, b):find("_")) + a - 1
    print(line:sub(1, curr-1) .. ((curr < px) and "/" or (curr == px) and "|" or "\\") .. line:sub(curr+1))
  end
end
