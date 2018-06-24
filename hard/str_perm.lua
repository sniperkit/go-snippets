function permgen(a, n)
  n = n or #a
  if n <= 1 then
    coroutine.yield(a)
  else
    for i = 1, n do
      a[n], a[i] = a[i], a[n]
      permgen(a, n-1)
      a[n], a[i] = a[i], a[n]
    end
  end
end

function permutations(a)
  local co = coroutine.create(function() permgen(a) end)
  return function()
    local code, res = coroutine.resume(co)
    return res
  end
end

for line in io.lines(arg[1]) do
  l, a = {}, {}
  for i in line:gmatch(".") do
    l[#l+1] = i
  end
  for i in permutations(l) do
    a[#a+1] = table.concat(i)
  end
  table.sort(a)
  print(table.concat(a, ","))
end
