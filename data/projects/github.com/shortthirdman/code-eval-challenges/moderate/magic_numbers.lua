function div(a, b) return math.floor(a/b) end

function isMagic(a)
  local d, r, ns = {}, 0, {}
  while a > 0 do
    r = a % 10
    if r == 0 or d[r] then return false end
    d[r] = true
    ns[#ns+1] = r
    a = div(a, 10)
  end
  d, r = {}, 1
  for i = 1, #ns do
    r = (r + ns[#ns - r + 1]) % #ns
    r = (r == 0) and #ns or r
    if d[r] then return false end
    d[r] = true
  end
  return r == 1
end

local maxb, magic = 0, {}
for line in io.lines(arg[1]) do
  local a, b = line:match("(%d+) (%d+)")
  a, b = tonumber(a), tonumber(b)
  while b > maxb do
    maxb = maxb + 1
    if isMagic(maxb) then
      magic[#magic+1] = maxb
    end
  end
  local r = {}
  for i = 1, #magic do
    if magic[i] > b then
      break
    elseif magic[i] >= a then
      r[#r+1] = magic[i]
    end
  end
  print((#r > 0) and table.concat(r, " ") or -1)
end
