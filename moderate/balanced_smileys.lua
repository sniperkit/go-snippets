function isBalanced(s, c)
  while true do
    if c < 0 then
      return false
    elseif s == "" then
      return c == 0
    end
    local first, last = s:sub(1,1), s:sub(-1)
    if first:find("[%a ]") then
      s = s:gsub("^[%a ]+", "")
    elseif last:find("[%a :]") then
      s = s:gsub("[%a :]+$", "")
    elseif first == "(" then
      c, s = c+1, s:sub(2)
    elseif first == ")" then
      c, s = c-1, s:sub(2)
    elseif first == ":" then
      if s:sub(2,2) == "(" then
        return isBalanced(s:sub(3), c) or isBalanced(s:sub(3), c+1)
      elseif s:sub(2,2) == ")" then
        return isBalanced(s:sub(3), c) or isBalanced(s:sub(3), c-1)
      else
        s = s:sub(2)
      end
    else
      return false
    end
  end
end

for line in io.lines(arg[1]) do
  print(isBalanced(line, 0) and "YES" or "NO")
end
