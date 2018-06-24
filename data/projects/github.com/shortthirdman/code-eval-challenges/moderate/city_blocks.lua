for line in io.lines(arg[1]) do
  local p, q = line:match("%(([%d,]+)%) %(([%d,]+)%)")
  local st, av = {}, {}
  for i in p:gmatch("%d+") do st[#st + 1] = tonumber(i) end
  for i in q:gmatch("%d+") do av[#av + 1] = tonumber(i) end
  if #st < 2 or #av < 2 then
    print(0)
  else
    if st[1] ~= 0 then
      for i = #st, 1, -1 do st[i] = st[i] - st[1] end
    end
    if av[1] ~= 0 then
      for i = #av, 1, -1 do av[i] = av[i] - av[1] end
    end
    local intersec, stl, avl = 0, st[#st], av[#av]
    for i = 1, #av do av[i] = av[i] * stl / avl end
    local i, j = 0, 0
    while i < #st and j < #av do
      if st[i] == av[j] then
        intersec = intersec + 1
        i = i + 1
        j = j + 1
      elseif st[i] < av[j] then
        i = i + 1
      else
        j = j + 1
      end
    end
    print(#st + #av - intersec - 1)
  end
end
