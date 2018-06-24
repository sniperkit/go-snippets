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
  a, f = {}, false
  for i in line:gmatch("%d+") do
    a[#a+1] = tonumber(i)
  end
  for i in permutations(a) do
    for o1 = 1, 3 do
      r1 = (o1 == 1) and i[1] + i[2] or (o1 == 2) and i[1] - i[2] or i[1] * i[2]
      for o2 = 1, 3 do
        r2 = (o2 == 1) and r1 + i[3] or (o2 == 2) and r1 - i[3] or r1 * i[3]
        for o3 = 1, 3 do
          r3 = (o3 == 1) and r2 + i[4] or (o3 == 2) and r2 - i[4] or r2 * i[4]
          for o4 = 1, 3 do
            r4 = (o4 == 1) and r3 + i[5] or (o4 == 2) and r3 - i[5] or r3 * i[5]
            if r4 == 42 then
              f = true
              break
            end
          end
          if f then break end
        end
        if f then break end
      end
      if f then break end
    end
    if f then break end
  end
  print(f and "YES" or "NO")
end
