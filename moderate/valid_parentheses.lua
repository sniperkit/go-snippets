for line in io.lines(arg[1]) do
  local valid, stack = true, {}
  for i in line:gmatch(".") do
    if i == ")" then
      if #stack > 0 and stack[#stack] == "(" then
        table.remove(stack)
      else
        valid = false
        break
      end
    elseif i == "]" then
      if #stack > 0 and stack[#stack] == "[" then
        table.remove(stack)
      else
        valid = false
        break
      end
    elseif i == "}" then
      if #stack > 0 and stack[#stack] == "{" then
        table.remove(stack)
      else
        valid = false
        break
      end
    else
      stack[#stack+1] = i
    end
  end
  print((valid and #stack == 0) and "True" or "False")
end
