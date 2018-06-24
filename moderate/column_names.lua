function div(a, b) return math.floor(a/b) end

local o = string.byte('A')
for a in io.lines(arg[1]) do
  a = a - 1
  local d = string.char(a % 26 + o)
  if a >= 26 then
    a = div(a, 26) - 1
    d = string.char(a % 26 + o) .. d
    if a >= 26 then
      a = div(a, 26) - 1
      d = string.char(a % 26 + o) .. d
    end
  end
  print(d)
end
