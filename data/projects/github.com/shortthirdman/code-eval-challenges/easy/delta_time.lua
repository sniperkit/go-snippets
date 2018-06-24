function div(a, b) return math.floor(a/b) end

for line in io.lines(arg[1]) do
  local a = {}
  for i in line:gmatch("%d%d") do a[#a+1] = i end
  local t = math.abs((a[1]-a[4])*3600 + (a[2]-a[5])*60 + a[3]-a[6])
  io.write(string.format("%02d:%02d:%02d\n", div(t, 3600), div(t, 60)%60, t%60))
end
