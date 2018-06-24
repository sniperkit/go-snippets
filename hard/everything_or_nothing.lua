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
function band(a,b) return (a + b - xor(a, b))/2 end

function check(a, b)
  return band(a, b) == 0
end

read, write, grant = 4, 2, 1

for line in io.lines(arg[1]) do
  local uperm, v = {7, 3, 0, 6, 2, 4, 5, 1, 5, 3, 7, 1, 6, 0, 2, 4, 2, 6}, true
  for i in line:gmatch("%S+") do
    local u, f, a = i:match("user_(%d)=>file_(%d)=>(%l)")
    local x = (tonumber(u) - 1) * 3 + tonumber(f)
    if a == 'r' and check(uperm[x], read) then
      v = false
      break
    elseif a == 'w' and check(uperm[x], write) then
      v = false
      break
    elseif a == 'g' then
      if check(uperm[x], grant) then
        v = false
        break
      end
      a, u = i:match("grant=>(%l)%l+=>user_(%d)")
      x = (tonumber(u) - 1) * 3 + tonumber(f)
      if a == 'r' and check(uperm[x], read) then
        uperm[x] = uperm[x] + read
      elseif a == 'w' and check(uperm[x], write) then
        uperm[x] = uperm[x] + write
      elseif a == 'g' and check(uperm[x], grant) then
        uperm[x] = uperm[x] + grant
      end
    end
  end
  print(v and "True" or "False")
end
