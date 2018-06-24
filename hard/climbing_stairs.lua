plen = 13
function addp(x, y, c, f)
  local ret = tostring(tonumber(x) + tonumber(y) + c)
  if f then return ret, 0 end
  if #ret < plen then
    return string.rep("0", plen - #ret) .. ret, 0
  elseif #ret > plen then
    return ret:sub(2), 1
  end
  return ret, 0
end

function add(x, y)
  if #x > #y then a, b = x, y else a, b = y, x end
  local r, c = "", 0
  local i, j = #a, #b
  while i > 0 do
    i, j = i - plen, j - plen
    local s, t, f = "", "", false
    if i <= 0 then
      s, f = a:sub(1, i+plen), true
    else
      s = a:sub(i+1, i+plen)
    end
    t = (j <= -plen) and "0" or (j <= 0) and b:sub(1, j+plen) or b:sub(j+1, j+plen)
    sum, c = addp(s, t, c, f)
    r = sum .. r
  end
  return r
end

for line in io.lines(arg[1]) do
  if #line > 0 then
    local a = tonumber(line)
    if a < 2 then
      print(a)
    else
      local b, c = "1", "1"
      while a > 1 do
        a, b, c = a-1, add(b, c), b
      end
      print(b)
    end
  end
end
