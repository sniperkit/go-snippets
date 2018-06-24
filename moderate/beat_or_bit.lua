function div(a, b) return math.floor(a/b) end

function bin2dec(s)
  b = 0
  for i = 1, #s do
    if s:sub(i, i) == '1' then
      b = b + 2 ^ (#s-i)
    end
  end
  return b
end

XOR_l = {{0,1}, {1,0}}
function xor(a, b)
  pow, c = 1, 0
  while a > 0 or b > 0 do
    c = c + (XOR_l[(a % 2)+1][(b % 2)+1] * pow)
    a, b, pow = div(a, 2), div(b, 2), pow * 2
  end
  return c
end

function g2d (a)
  a = xor(a, div(a, 16))
  a = xor(a, div(a, 4))
  a = xor(a, div(a, 2))
  return a
end

for line in io.lines(arg[1]) do
  r = {}
  for i in line:gmatch("%d+") do
    r[#r + 1] = g2d(bin2dec(i))
  end
  print(table.concat(r, " | "))
end
