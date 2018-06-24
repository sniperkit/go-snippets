track, cars = {}, {}
for line in io.lines(arg[1]) do
  if #track == 0 then
    for i in line:gmatch("%S+ %d+") do
      local len, angle = i:match("(%S+) (%d+)")
      track[#track + 1] = {tonumber(len), tonumber(angle)}
    end
  else
    local id, vmax, accel, brake = line:match("(%d+) (%d+) (%S+) (%S+)")
    local laptime, vinit = 0, 0
    for _, i in ipairs(track) do
      local len, angle = i[1], i[2]
      local away = accel * (1 - (vinit / vmax)) * (vinit + vmax) / (2 * 3600)
      local vend = vmax * (1 - (angle / 180))
      local bway = brake * (1 - (vend / vmax)) * (vend + vmax) / (2 * 3600)
      laptime = laptime + 3600 * (2*away/(vinit+vmax) + (len-away-bway)/vmax + 2*bway/(vend+vmax))
      vinit = vend
    end
    cars[#cars + 1] = {id, laptime}
  end
end
table.sort(cars, function (a,b) return (a[2] < b[2]) end)
for _, i in ipairs(cars) do
  print(string.format("%s %.2f", i[1], i[2]))
end
