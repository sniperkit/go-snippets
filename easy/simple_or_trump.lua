v = {10, 2, 3, 4, 5, 6, 7, 8, 9, J = 11, Q = 12, K = 13, A = 14}
function rank (a, t)
  local r
  if tonumber(a:sub(1, 1)) ~= nil then
    r = v[tonumber(a:sub(1, 1))]
  else
    r = v[a:sub(1, 1)]
  end
  if a:sub(#a, #a) == t then
    r = r + 13
  end
  return r
end

function win (l, r, t)
  local lr, rr = rank(l, t), rank(r, t)
  if lr > rr then
    return l
  elseif lr < rr then
    return r
  end
  return l .. " " .. r
end

for line in io.lines(arg[1]) do
  local l, r, t = line:match("(%w+) (%w+) | (%u)")
  print(win(l, r, t))
end
