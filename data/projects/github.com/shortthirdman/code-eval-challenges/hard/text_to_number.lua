dn, dz = "negative", "zero"
d0 = {zero=0, one=1, two=2, three=3, four=4, five=5, six=6, seven=7,
      eight=8, nine=9}
d1 = {ten=10, eleven=11, twelve=12, thirteen=13, fourteen=14,
      fifteen=15, sixteen=16, seventeen=17, eighteen=18, nineteen=19}
d2 = {twenty=20, thirty=30, forty=40, fifty=50, sixty=60, seventy=70,
      eighty=80, ninety=90}
d3 = {"hundred", "thousand", "million"}

for line in io.lines(arg[1]) do
  local n, h, q, r, t = 1, 0, 0, 0, {}
  for i in line:gmatch("%S+") do
    if i == dn then
      n = -1
    else
      t[#t+1] = i
    end
  end
  for i = 1, #t do
    if t[i] == dz then
      q = 0
    elseif d0[t[i]] then
      h = d0[t[i]]
    elseif d1[t[i]] then
      q = q + d1[t[i]]
    elseif d2[t[i]] then
      q = q + d2[t[i]]
    elseif t[i] == d3[1] then
      q, h = h*100, 0
    elseif t[i] == d3[2] then
      r, h, q = r+(h+q)*1000, 0, 0
    elseif t[i] == d3[3] then
      r, h, q = (h+q)*1000000, 0, 0
    end
  end
  print(n * (r + q + h))
end
