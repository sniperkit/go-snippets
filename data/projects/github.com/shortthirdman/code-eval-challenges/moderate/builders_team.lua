function div(a, b) return math.floor(a/b) end

XOR_l = {{0,1}, {1,0}}
function xor(a, b)
   pow, c = 1, 0
   while a > 0 or b > 0 do
      c = c + (XOR_l[(a % 2)+1][(b % 2)+1] * pow)
      a, b, pow = div(a, 2), div(b, 2), pow * 2
   end
   return c
end
local ff = 2^32 - 1
function band(a,b) return (a + b - xor(a, b))/2 end
function bor (a,b) return ff - band(ff - a, ff - b) end

hor = { 0x42, 0x84, 0x108, 0x210,
	0x840, 0x1080, 0x2100, 0x4200,
	0x10800, 0x21000, 0x42000, 0x84000,
	0x210000, 0x420000, 0x840000, 0x1080000,
	0x1806, 0x300c, 0x6018,
	0x300c0, 0x60180, 0xc0300,
	0x601800, 0xc03000, 0x1806000,
	0x7000e, 0xe001c,
	0xe001c0, 0x1c00380,
	0x1e0001e }
ver = { 0x6, 0xc, 0x18, 0x30,
	0xc0, 0x180, 0x300, 0x600,
	0x1800, 0x3000, 0x6000, 0xc000,
	0x30000, 0x60000, 0xc0000, 0x180000,
	0x14a, 0x294, 0x528,
	0x2940, 0x5280, 0xa500,
	0x52800, 0xa5000, 0x14a000,
	0x4a52, 0x94a4,
	0x94a40, 0x129480,
	0x118c62 }

for line in io.lines(arg[1]) do
  local h, v, r = 0, 0, 0
  for i in line:gmatch("%d+ %d+") do
    local a, b = i:match("(%d+) (%d+)")
    a, b = tonumber(a), tonumber(b)
    if a > b then a, b = b, a end
    if a + 1 == b then
      h = bor(h, 2^a)
    else
      v = bor(v, 2^a)
    end
  end
  for i = 1, 30 do
    if band(h, hor[i]) == hor[i] and band(v, ver[i]) == ver[i] then
      r = r + 1
    end
  end
  print(r)
end
