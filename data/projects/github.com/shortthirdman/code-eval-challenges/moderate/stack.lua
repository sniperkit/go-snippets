function push(stack, int)
  stack[#stack+1] = int
  return
end

function pop(stack)
  local int = stack[#stack]
  stack[#stack] = nil
  return int
end

for line in io.lines(arg[1]) do
  local stack = {}
  for i in line:gmatch("%S+") do
    push(stack, tonumber(i))
  end
  local first = true
  while #stack > 0 do
    if first then first = false else io.write(" ") end
    io.write(pop(stack))
    pop(stack)
  end
  print()
end
