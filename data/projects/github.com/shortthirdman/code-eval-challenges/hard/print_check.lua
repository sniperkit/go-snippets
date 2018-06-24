function div(a, b) return math.floor(a/b) end

s = {[0]="", "One", "Two", "Three", "Four",
     "Five", "Six", "Seven", "Eight", "Nine",
     "Ten", "Eleven", "Twelve", "Thirteen", "Fourteen",
     "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen",
     "Twenty", "Thirty", "Forty", "Fifty",
     "Sixty", "Seventy", "Eighty", "Ninety",
     "Hundred", "Thousand", "Million"}

function wrd(a1, a2, a3)
  return (a1 + a2 + a3 > 0) and (((a1 > 0) and s[a1] .. s[28] or "") ..
                                 ((a2 > 1) and s[18 + a2] .. s[a3] or s[10*a2 + a3]))
end

for line in io.lines(arg[1]) do
  local b, c = tonumber(line), {}
  if b == 0 then
    io.write("Zero")
  else
    for i = 9, 1, -1 do
      c[i], b = b%10, div(b, 10)
    end
    local r = wrd(c[1], c[2], c[3])
    if r then io.write(r .. s[30]) end
    r = wrd(c[4], c[5], c[6])
    if r then io.write(r .. s[29]) end
    r = wrd(c[7], c[8], c[9])
    if r then io.write(r) end
  end
  io.write("Dollars\n")
end
