for line in io.lines(arg[1]) do
  local l, m = #line, ""
  if l > 0 then
    local c = 1
    for n = 1, (l-c+1)/2 do
      local f = false
      for i = c, l-n do
        subs = line:sub(i, i+n-1)
        if not subs:find("^%s+$") and string.gsub(line:sub(i+n, l), "%%", "%%"):find(subs) then
          m, c, f = subs, i, true
          break
        end
      end
      if f == false then break end
    end
  end
  print((m == "") and "NONE" or m)
end
