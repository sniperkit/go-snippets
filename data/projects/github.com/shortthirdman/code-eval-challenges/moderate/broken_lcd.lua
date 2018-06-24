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
function band(a,b) return (a + b - xor(a, b))/2 end

digits = { 252,  96, 218, 242, 102, 182, 190, 224, 254, 246 }

for line in io.lines(arg[1]) do
  sep, a, d = line:find(";"), {}, {}
  for i in line:sub(1, sep):gmatch("%d+") do
    a[#a + 1] = bin2dec(i)
  end
  dig = line:sub(sep + 1)
  d[1], j = digits[tonumber(dig:sub(1,1))+1], 2
  for i = 2, #dig do
    if dig:sub(i,i) == '.' then
      d[j-1] = d[j-1] + 1
    else
      d[j] = digits[tonumber(dig:sub(i,i))+1]
      j = j + 1
    end
  end
  for i = 0, 13 - #dig do
    f = true
    for j = 1, #dig - 1 do
      if band(a[i + j], d[j]) ~= d[j] then
        f = false
        break
      end
    end
    if f then
      break
    end
  end
  print(f and 1 or 0)
end
