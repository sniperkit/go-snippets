function tadd(s, t)
  for i = 1, #t do
    s[#s + 1] = t[i]
  end
  return s
end

function tele(s, t)
  local ret = {}
  if #s == 0 then
    ret[1] = t
  elseif s:sub(1, 1) == '2' then
    ret = tele(s:sub(2), t .. "a")
    ret = tadd(ret, tele(s:sub(2), t .. "b"))
    ret = tadd(ret, tele(s:sub(2), t .. "c"))
  elseif s:sub(1, 1) == '3' then
    ret = tele(s:sub(2), t .. "d")
    ret = tadd(ret, tele(s:sub(2), t .. "e"))
    ret = tadd(ret, tele(s:sub(2), t .. "f"))
  elseif s:sub(1, 1) == '4' then
    ret = tele(s:sub(2), t .. "g")
    ret = tadd(ret, tele(s:sub(2), t .. "h"))
    ret = tadd(ret, tele(s:sub(2), t .. "i"))
  elseif s:sub(1, 1) == '5' then
    ret = tele(s:sub(2), t .. "j")
    ret = tadd(ret, tele(s:sub(2), t .. "k"))
    ret = tadd(ret, tele(s:sub(2), t .. "l"))
  elseif s:sub(1, 1) == '6' then
    ret = tele(s:sub(2), t .. "m")
    ret = tadd(ret, tele(s:sub(2), t .. "n"))
    ret = tadd(ret, tele(s:sub(2), t .. "o"))
  elseif s:sub(1, 1) == '7' then
    ret = tele(s:sub(2), t .. "p")
    ret = tadd(ret, tele(s:sub(2), t .. "q"))
    ret = tadd(ret, tele(s:sub(2), t .. "r"))
    ret = tadd(ret, tele(s:sub(2), t .. "s"))
  elseif s:sub(1, 1) == '8' then
    ret = tele(s:sub(2), t .. "t")
    ret = tadd(ret, tele(s:sub(2), t .. "u"))
    ret = tadd(ret, tele(s:sub(2), t .. "v"))
  elseif s:sub(1, 1) == '9' then
    ret = tele(s:sub(2), t .. "w")
    ret = tadd(ret, tele(s:sub(2), t .. "x"))
    ret = tadd(ret, tele(s:sub(2), t .. "y"))
    ret = tadd(ret, tele(s:sub(2), t .. "z"))
  else
    return tele(s:sub(2), t .. s:sub(1, 1))
  end
  return ret
end

for line in io.lines(arg[1]) do
  print(table.concat(tele(line, ""), ","))
end
